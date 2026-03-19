package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"log/slog"

	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/common/utils"
	"github.com/baowk/dilu-core/common/utils/ips"
	"github.com/baowk/dilu-core/common/utils/text"
	"github.com/baowk/dilu-core/config"
	"github.com/baowk/dilu-core/core/cache"
	"github.com/baowk/dilu-core/core/logger"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var app *Application

// Application 核心应用结构，替代全局变量
type Application struct {
	config     config.Config
	logger     *slog.Logger
	cache      cache.ICache
	redisLock  *redislock.Client
	engine     http.Handler
	databases  map[string]*gorm.DB
	started    chan struct{}
	toClose    chan struct{}
	mu         sync.RWMutex
	engineOnce sync.Once
	dbInitFlag bool
	health     *HealthService
	monitor    *Monitor
	registry   ServiceRegistry
}

// Init 初始化应用
func Init(cfg config.Config) error {
	app = &Application{
		config:    cfg,
		databases: make(map[string]*gorm.DB),
		started:   make(chan struct{}, 1),
		toClose:   make(chan struct{}, 1),
	}

	// 初始化日志
	app.logger = logger.InitLogger(*cfg.GetLogCfg())

	utils.Setup(cfg.GetServerCfg().GetNode())

	// 初始化缓存
	c, err := cache.New(*cfg.GetCacheCfg())
	if err != nil {
		return fmt.Errorf("cache initialization failed: %w", err)
	}
	app.cache = c
	if app.cache.Type() == "redis" {
		if r, ok := app.cache.(*cache.RedisCache); ok {
			app.redisLock = redislock.New(r.GetClient())
		}
	}

	// 初始化数据库
	if err := app.dbInit(); err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}

	app.dbInitFlag = true

	// 初始化健康检查和监控
	app.health = NewHealthService(app.logger)
	app.monitor = NewMonitor(app)
	app.health.RegisterChecker(NewDatabaseHealthChecker(app))
	app.health.RegisterChecker(NewCacheHealthChecker(app))

	return nil
}

func GetApp() *Application {
	return app
}

// GetConfig 获取配置
func (app *Application) GetConfig() config.Config {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.config
}

// SetConfig 设置配置
func (app *Application) SetConfig(cfg config.Config) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.config = cfg
}

// GetLogger 获取日志实例
func (app *Application) GetLogger() *slog.Logger {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.logger
}

// GetCache 获取缓存实例
func (app *Application) GetCache() cache.ICache {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.cache
}

func (app *Application) Db(name string) *gorm.DB {
	app.mu.RLock()
	defer app.mu.RUnlock()
	if name == "" {
		name = consts.DB_DEF
	}
	return app.databases[name]
}

// GetEngine 获取HTTP引擎
func (app *Application) GetEngine() http.Handler {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.engine
}

// SetEngine 设置HTTP引擎
func (app *Application) SetEngine(engine http.Handler) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.engine = engine
}

// GetGinEngine 获取Gin引擎
func (app *Application) GetGinEngine() *gin.Engine {
	app.engineOnce.Do(func() {
		cfg := app.GetConfig()
		if cfg.GetServerCfg().Mode == ModeProd.String() {
			gin.SetMode(gin.ReleaseMode)
		}
		engine := gin.New()
		app.mu.Lock()
		app.engine = engine
		app.mu.Unlock()

		// 自动注册健康检查路由（直接传 engine，避免 GetGinEngine 死锁）
		app.registerHealthRoutesOn(engine)
	})

	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.engine.(*gin.Engine)
}

// Run 运行应用
func (app *Application) Run() error {
	cfg := app.GetConfig()
	addr := fmt.Sprintf("%s:%d", cfg.GetServerCfg().GetHost(), cfg.GetServerCfg().GetPort())

	// 服务启动参数
	srv := &http.Server{
		Addr:           addr,
		Handler:        app.GetEngine(),
		ReadTimeout:    time.Duration(cfg.GetServerCfg().GetReadTimeout()) * time.Second,
		WriteTimeout:   time.Duration(cfg.GetServerCfg().GetWriteTimeout()) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	// 打印启动信息
	fmt.Println(text.Green(`Dilu github:`) + text.Blue(`https://github.com/baowk/dilu`))
	fmt.Println(text.Green(app.config.GetServerCfg().Name) + fmt.Sprintf(" %d ", app.config.GetServerCfg().GetNode()) + text.Green("Server started ,Listen: ") + text.Red("[ "+addr+" ]"))
	fmt.Println(text.Yellow("Dilu Go Go Go ~ ~ ~ "))

	if cfg.GetServerCfg().Mode != ModeProd.String() {
		fmt.Println(text.Blue(fmt.Sprintf("Swagger: http://localhost:%d/swagger/index.html", cfg.GetServerCfg().GetPort())))
		if ip := ips.GetLocalHost(); ip != "" {
			fmt.Println(text.Blue(fmt.Sprintf("Swagger: http://%s:%d/swagger/index.html", ip, cfg.GetServerCfg().GetPort())))
		}
	}

	// 自动注册到服务注册中心
	if app.registry != nil {
		healthURL := fmt.Sprintf("http://%s:%d/health", cfg.GetServerCfg().GetHost(), cfg.GetServerCfg().GetPort())
		if cfg.GetServerCfg().GetHost() == "0.0.0.0" {
			if localIP := ips.GetLocalHost(); localIP != "" {
				healthURL = fmt.Sprintf("http://%s:%d/health", localIP, cfg.GetServerCfg().GetPort())
			}
		}
		var tags []string
		if ac, ok := cfg.(*config.AppConfig); ok {
			tags = ac.Rd.Tags
		}
		if err := app.registry.Register(
			cfg.GetServerCfg().Name,
			cfg.GetServerCfg().GetHost(),
			cfg.GetServerCfg().GetPort(),
			healthURL,
			tags,
		); err != nil {
			app.GetLogger().Error("服务注册失败", "error", err)
		} else {
			app.GetLogger().Info("服务已注册到注册中心")
		}
	}

	app.GetLogger().Debug("服务初始化完毕")
	close(app.started)
	app.GetLogger().Debug("给信号量Started")

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务器信号
	close(app.toClose)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.GetLogger().Info("Shutdown Server ...", "time", time.Now())

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// 从注册中心注销
	if app.registry != nil {
		app.registry.Deregister()
		app.GetLogger().Info("服务已从注册中心注销")
	}

	app.GetLogger().Info("Server exiting")
	time.Sleep(time.Second * time.Duration(cfg.GetServerCfg().GetCloseWait()))
	return nil
}

// CacheRedis 获取Redis客户端
func (app *Application) CacheRedis() (redis.UniversalClient, error) {
	return cache.GetRedisClient(app.GetCache())
}

// WaitForStart 等待应用启动完成
func (app *Application) WaitForStart() <-chan struct{} {
	return app.started
}

// WaitForClose 等待应用关闭信号
func (app *Application) WaitForClose() <-chan struct{} {
	return app.toClose
}

func GetCache() cache.ICache {
	return app.cache
}

// 获取默认的（master）db
func DB() *gorm.DB {
	return app.Db(consts.DB_DEF)
}

func GetGinEngine() *gin.Engine {
	return app.GetGinEngine()
}

func GetRedisLock() *redislock.Client {
	return app.redisLock
}

// GetHealthService 获取健康检查服务
func (app *Application) GetHealthService() *HealthService {
	return app.health
}

// GetMonitor 获取监控器
func (app *Application) GetMonitor() *Monitor {
	return app.monitor
}

// registerHealthRoutesOn 在指定的 Gin 引擎上注册健康检查和监控路由（内部使用，避免死锁）
func (app *Application) registerHealthRoutesOn(engine *gin.Engine) {
	// 健康检查端点 - 供 dilu-rd 心跳检测使用
	engine.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		results := app.health.CheckAll(ctx)
		status := app.health.GetOverallStatus(results)

		httpCode := http.StatusOK
		if status == "unhealthy" {
			httpCode = http.StatusServiceUnavailable
		}

		c.JSON(httpCode, gin.H{
			"status":  status,
			"checks":  results,
			"service": app.config.GetServerCfg().Name,
		})
	})

	// 监控指标端点
	engine.GET("/metrics", func(c *gin.Context) {
		metrics := app.monitor.CollectMetrics()
		c.JSON(http.StatusOK, metrics)
	})

	// 路由导出端点 - 供 dilu-gateway 自动发现
	engine.GET("/routes", func(c *gin.Context) {
		type RouteInfo struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}
		routes := engine.Routes()
		result := make([]RouteInfo, 0, len(routes))
		for _, r := range routes {
			// 排除内部管理端点
			if r.Path == "/health" || r.Path == "/metrics" || r.Path == "/routes" {
				continue
			}
			result = append(result, RouteInfo{
				Method: r.Method,
				Path:   r.Path,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"service": app.config.GetServerCfg().Name,
			"routes":  result,
			"total":   len(result),
		})
	})
}

// RegisterHealthRoutes 公开方法，向当前 Gin 引擎注册健康检查路由
func (app *Application) RegisterHealthRoutes() {
	app.registerHealthRoutesOn(app.GetGinEngine())
}

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
	dbInitFlag bool
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
	app.cache = cache.New(*cfg.GetCacheCfg())
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
	cfg := app.GetConfig()
	if cfg.GetServerCfg().Mode == ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}

	app.mu.RLock()
	currentEngine := app.engine
	app.mu.RUnlock()

	if currentEngine == nil {
		app.mu.Lock()
		if app.engine == nil {
			app.engine = gin.New()
		}
		currentEngine = app.engine
		app.mu.Unlock()
	}

	switch engine := currentEngine.(type) {
	case *gin.Engine:
		return engine
	default:
		panic("not support other engine")
	}
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

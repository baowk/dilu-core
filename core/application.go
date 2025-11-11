package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"time"

	"os"

	"github.com/baowk/dilu-core/common/utils/ips"
	"github.com/baowk/dilu-core/common/utils/text"
	"github.com/baowk/dilu-core/config"
	"github.com/baowk/dilu-core/core/cache"
	"github.com/baowk/dilu-core/core/locker"
	"github.com/baowk/dilu-core/core/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"log/slog"
)

var (
	Cfg        config.AppCfg
	Log        *slog.Logger //*zap.Logger
	Cache      cache.ICache
	dbInitFlag bool // 数据库是否初始化
	engine     http.Handler
	dbs        = make(map[string]*gorm.DB, 0)
	RedisLock  *locker.Redis
	Started    = make(chan byte, 1)
	ToClose    = make(chan byte, 1)
	//lock      sync.RWMutex
)

func GetEngine() http.Handler {
	return engine
}

func SetEngine(aEngine http.Handler) {
	engine = aEngine
}

func GetGinEngine() *gin.Engine {
	if Cfg.Server.Mode == ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	var r *gin.Engine
	// lock.RLock()
	// defer lock.RUnlock()
	if engine == nil {
		engine = gin.New()
	}
	switch engine := engine.(type) {
	case *gin.Engine:
		r = engine
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}
	return r
}

func Init() {
	Log = logger.InitLogger(Cfg.Logger)
	Cache = cache.New(Cfg.Cache)
	if Cache.Type() == "redis" {
		r := Cache.(*cache.RedisCache)
		RedisLock = locker.NewRedis(r.GetClient())
	}
	dbInit()
}

func Run() {
	addr := fmt.Sprintf("%s:%d", Cfg.Server.GetHost(), Cfg.Server.GetPort())

	//服务启动参数
	srv := &http.Server{
		Addr:           addr,
		Handler:        GetEngine(),
		ReadTimeout:    time.Duration(Cfg.Server.GetReadTimeout()) * time.Second,
		WriteTimeout:   time.Duration(Cfg.Server.GetWriteTimeout()) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	fmt.Println(text.Green(`Dilu github:`) + text.Blue(`https://github.com/baowk/dilu`))
	fmt.Println(text.Green("Dilu Server started ,Listen: ") + text.Red("[ "+addr+" ]"))
	fmt.Println(text.Yellow("Dilu Go Go Go ~ ~ ~ "))

	if Cfg.Server.Mode != ModeProd.String() {
		//fmt.Printf("Swagger %s %s start\r\n", docs.SwaggerInfo.Title, docs.SwaggerInfo.Version)
		fmt.Println(text.Blue(fmt.Sprintf("Swagger: http://localhost:%d/swagger/index.html", Cfg.Server.Port)))
		ip := ips.GetLocalHost()
		if ip != "" {
			fmt.Println(text.Blue(fmt.Sprintf("Swagger: http://%s:%d/swagger/index.html", ip, Cfg.Server.Port)))
		}
	}

	slog.Debug("服务初始化完毕")
	Started <- 1
	slog.Debug("给信号量Started")
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	//关闭服务器信号
	ToClose <- 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slog.Info("Shutdown Server ...", "time", time.Now())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	slog.Info("Server exiting")
	time.Sleep(time.Second * time.Duration(Cfg.Server.GetCloseWait()))
}

// 此方法仅用于配置了redis缓存
func CacheRedis() (redis.UniversalClient, error) {
	return cache.GetRedisClient(Cache)
}

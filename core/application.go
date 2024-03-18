package core

import (
	"context"
	"fmt"
	"github.com/baowk/dilu-core/core/mongodb"
	"log"
	"net/http"
	"os/signal"
	"time"

	"os"
	"sync"

	"github.com/baowk/dilu-core/common/utils/files"
	"github.com/baowk/dilu-core/common/utils/ips"
	"github.com/baowk/dilu-core/common/utils/text"
	"github.com/baowk/dilu-core/config"
	"github.com/baowk/dilu-core/core/cache"
	"github.com/baowk/dilu-core/core/locker"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

var (
	Cfg       config.AppCfg
	Log       *zap.Logger
	Cache     cache.ICache
	lock      sync.RWMutex
	engine    http.Handler
	dbs       = make(map[string]*gorm.DB, 0)
	RedisLock *locker.Redis
)

func GetEngine() http.Handler {
	return engine
}

func SetEngine(aEngine http.Handler) {
	engine = aEngine
}

func GetGinEngine() *gin.Engine {
	var r *gin.Engine
	lock.RLock()
	defer lock.RUnlock()
	if engine == nil {
		engine = gin.New()
	}
	switch engine.(type) {
	case *gin.Engine:
		r = engine.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}
	return r
}

func Init() {
	logInit()
	Cache = cache.New(Cfg.Cache)
	if Cache.Type() == "redis" {
		r := Cache.(*cache.RedisCache)
		RedisLock = locker.NewRedis(r.GetClient())
	}
	dbInit()
	// 判断是否启用mongo
	if Cfg.Mongodb.Open {
		mongodb.MongoInit(Cfg.Mongodb)
	}
}

// RunMode gin运行模式
func RunMode() {
	if Cfg.Server.Mode == ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
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

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Printf("%s Shutdown Server ... \r\n", time.Now())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func logInit() {
	//初始化日志
	Log = zapInit()
	zap.ReplaceGlobals(Log)
}

// Zap 获取 zap.Logger
func zapInit() (logger *zap.Logger) {
	if ok, _ := files.PathExists(Cfg.Logger.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", Cfg.Logger.Director)
		_ = os.MkdirAll(Cfg.Logger.Director, os.ModePerm)
	}

	cores := Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if Cfg.Logger.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

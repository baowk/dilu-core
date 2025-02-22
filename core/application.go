package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"time"

	"os"
	"sync"

	"github.com/baowk/dilu-core/common/utils/ips"
	"github.com/baowk/dilu-core/common/utils/text"
	"github.com/baowk/dilu-core/config"
	"github.com/baowk/dilu-core/core/cache"
	"github.com/baowk/dilu-core/core/locker"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/redis/go-redis/v9"

	// "go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
	"log/slog"

	"gorm.io/gorm"
)

var (
	Cfg       config.AppCfg
	Log       *slog.Logger //*zap.Logger
	Cache     cache.ICache
	lock      sync.RWMutex
	engine    http.Handler
	dbs       = make(map[string]*gorm.DB, 0)
	RedisLock *locker.Redis
	Started   = make(chan byte, 1)
	ToClose   = make(chan byte, 1)
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
	logWrite := logInit()
	Cache = cache.New(Cfg.Cache)
	if Cache.Type() == "redis" {
		r := Cache.(*cache.RedisCache)
		RedisLock = locker.NewRedis(r.GetClient())
	}
	dbInit(logWrite)
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

func logInit() io.Writer {
	//初始化日志
	// Log = zapInit()
	// zap.ReplaceGlobals(Log)
	opts := slog.HandlerOptions{
		AddSource: Cfg.Logger.ShowLine,
		Level:     slog.LevelDebug,
	}
	var logW io.Writer

	if Cfg.Logger.LogInConsole {
		logW = os.Stdout
	} else {
		logW = &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: Cfg.Logger.Director + "/dilu.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: Cfg.Logger.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: Cfg.Logger.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: Cfg.Logger.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}
	}

	if Cfg.Logger.Level == "error" {
		opts.Level = slog.LevelError
	} else if Cfg.Logger.Level == "info" {
		opts.Level = slog.LevelInfo
	} else if Cfg.Logger.Level == "warn" {
		opts.Level = slog.LevelWarn
	}
	if Cfg.Logger.Format == "json" {
		Log = slog.New(slog.NewJSONHandler(logW, &opts))
	} else {
		Log = slog.New(slog.NewTextHandler(logW, &opts))
	}
	slog.SetDefault(Log)
	return logW
}

// 此方法仅用于配置了redis缓存
func CacheRedis() (redis.UniversalClient, error) {
	return cache.GetRedisClient(Cache)
}

// Zap 获取 zap.Logger
// func zapInit() (logger *zap.Logger) {
// 	if ok, _ := files.PathExists(Cfg.Logger.Director); !ok { // 判断是否有Director文件夹
// 		fmt.Printf("create %v directory\n", Cfg.Logger.Director)
// 		_ = os.MkdirAll(Cfg.Logger.Director, os.ModePerm)
// 	}

// 	cores := Zap.GetZapCores()
// 	logger = zap.New(zapcore.NewTee(cores...))

// 	if Cfg.Logger.ShowLine {
// 		logger = logger.WithOptions(zap.AddCaller())
// 	}
// 	return logger
// }

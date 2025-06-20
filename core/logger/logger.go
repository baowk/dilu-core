package logger

import (
	"log/slog"
	"os"

	"github.com/baowk/dilu-core/config"

	"github.com/natefinch/lumberjack"
)

func InitLogger(cfg config.AppCfg) {
	opts := slog.HandlerOptions{
		AddSource: cfg.Logger.ShowLine,
		Level:     slog.LevelDebug,
	}
	switch cfg.Logger.Level {
	case "error":
		opts.Level = slog.LevelError
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "debug":
		opts.Level = slog.LevelDebug
	}

	// 配置 LevelHandler
	levelHandler := NewLevelHandler()
	if opts.Level.Level() == slog.LevelDebug {

		// 始终创建文件写入器
		debugFile := &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: cfg.Logger.Director + "/debug.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: cfg.Logger.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: cfg.Logger.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: cfg.Logger.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}
		debugHandler := slog.NewJSONHandler(debugFile, &opts)
		levelHandler.AddHandler(slog.LevelDebug, debugHandler)
	}
	if opts.Level.Level() <= slog.LevelInfo {

		// 始终创建文件写入器
		infoFile := &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: cfg.Logger.Director + "/info.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: cfg.Logger.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: cfg.Logger.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: cfg.Logger.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}
		infoHandler := slog.NewJSONHandler(infoFile, &opts)
		levelHandler.AddHandler(slog.LevelInfo, infoHandler)
	}
	if opts.Level.Level() <= slog.LevelWarn {

		// 始终创建文件写入器
		warnFile := &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: cfg.Logger.Director + "/warn.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: cfg.Logger.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: cfg.Logger.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: cfg.Logger.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}
		warnHandler := slog.NewJSONHandler(warnFile, &opts)
		levelHandler.AddHandler(slog.LevelWarn, warnHandler)
	}
	if opts.Level.Level() <= slog.LevelError {
		// 始终创建文件写入器
		errorFile := &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: cfg.Logger.Director + "/error.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: cfg.Logger.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: cfg.Logger.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: cfg.Logger.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}
		errorHandler := slog.NewJSONHandler(errorFile, &opts)
		levelHandler.AddHandler(slog.LevelError, errorHandler)
	}

	if cfg.Logger.LogInConsole {
		// 同时输出到文件和控制台
		consoleHandler := slog.NewTextHandler(os.Stdout, &opts)
		if opts.Level.Level() == slog.LevelDebug {
			levelHandler.AddHandler(slog.LevelDebug, consoleHandler)
		}
		if opts.Level.Level() <= slog.LevelInfo {

			levelHandler.AddHandler(slog.LevelInfo, consoleHandler)
		}
		if opts.Level.Level() <= slog.LevelWarn {
			levelHandler.AddHandler(slog.LevelWarn, consoleHandler)
		}
		if opts.Level.Level() <= slog.LevelError {
			levelHandler.AddHandler(slog.LevelError, consoleHandler)
		}
	}

	// 设置默认 logger
	slog.SetDefault(slog.New(levelHandler))
}

package logger

import (
	"log/slog"
	"os"

	"github.com/baowk/dilu-core/config"

	"github.com/natefinch/lumberjack"
)

func InitLogger(logC config.LogCfg) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: logC.ShowLine,
		Level:     slog.LevelDebug,
	}
	switch logC.Level {
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
	if logC.OutputMode == "single" { //单个日志文件
		// 始终创建文件写入器
		diluFile := &lumberjack.Logger{
			// 日志文件名，归档日志也会保存在对应目录下
			// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
			// <processname>-lumberjack.log
			Filename: logC.Director + "/dilu.log",

			// backup的日志是否使用本地时间戳，默认使用UTC时间
			LocalTime: true,
			// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
			MaxSize: logC.GetMaxSize(),
			// 旧日志保存的最大天数，默认保存所有旧日志文件
			MaxAge: logC.GetMaxAge(),
			// 旧日志保存的最大数量，默认保存所有旧日志文件
			MaxBackups: logC.GetMaxBackups(),
			// 对backup的日志是否进行压缩，默认不压缩
			Compress: true,
		}

		//diluFile.Filename = logC.Director + "/dilu.log"

		var diluHandler slog.Handler
		if logC.Format == "json" {
			diluHandler = slog.NewJSONHandler(diluFile, &opts)
		} else {
			diluHandler = slog.NewTextHandler(diluFile, &opts)
		}

		if opts.Level.Level() == slog.LevelDebug {
			levelHandler.AddHandler(slog.LevelDebug, diluHandler)
		}
		if opts.Level.Level() <= slog.LevelInfo {
			levelHandler.AddHandler(slog.LevelInfo, diluHandler)
		}
		if opts.Level.Level() <= slog.LevelWarn {
			levelHandler.AddHandler(slog.LevelWarn, diluHandler)
		}
		if opts.Level.Level() <= slog.LevelError {
			levelHandler.AddHandler(slog.LevelError, diluHandler)
		}

	} else {
		if opts.Level.Level() == slog.LevelDebug {

			// 始终创建文件写入器
			debugFile := &lumberjack.Logger{
				// 日志文件名，归档日志也会保存在对应目录下
				// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
				// <processname>-lumberjack.log
				//Filename: logC.Director + "/debug.log",

				// backup的日志是否使用本地时间戳，默认使用UTC时间
				LocalTime: true,
				// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
				MaxSize: logC.GetMaxSize(),
				// 旧日志保存的最大天数，默认保存所有旧日志文件
				MaxAge: logC.GetMaxAge(),
				// 旧日志保存的最大数量，默认保存所有旧日志文件
				MaxBackups: logC.GetMaxBackups(),
				// 对backup的日志是否进行压缩，默认不压缩
				Compress: true,
			}

			// if logC.DebugFile != "" {
			// 	debugFile.Filename = logC.Director + "/" + logC.DebugFile
			// } else {
			debugFile.Filename = logC.Director + "/debug.log"
			//}

			if logC.Format == "json" {
				debugHandler := slog.NewJSONHandler(debugFile, &opts)
				levelHandler.AddHandler(slog.LevelDebug, debugHandler)
			} else {
				debugHandler := slog.NewTextHandler(debugFile, &opts)
				levelHandler.AddHandler(slog.LevelDebug, debugHandler)
			}
			// debugHandler := slog.NewJSONHandler(debugFile, &opts)
			// levelHandler.AddHandler(slog.LevelDebug, debugHandler)
		}
		if opts.Level.Level() <= slog.LevelInfo {

			// 始终创建文件写入器
			infoFile := &lumberjack.Logger{
				// 日志文件名，归档日志也会保存在对应目录下
				// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
				// <processname>-lumberjack.log
				//Filename: logC.Director + "/info.log",

				// backup的日志是否使用本地时间戳，默认使用UTC时间
				LocalTime: true,
				// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
				MaxSize: logC.GetMaxSize(),
				// 旧日志保存的最大天数，默认保存所有旧日志文件
				MaxAge: logC.GetMaxAge(),
				// 旧日志保存的最大数量，默认保存所有旧日志文件
				MaxBackups: logC.GetMaxBackups(),
				// 对backup的日志是否进行压缩，默认不压缩
				Compress: true,
			}
			// if logC.InfoFile != "" {
			// 	infoFile.Filename = logC.Director + "/" + logC.InfoFile
			// } else {
			infoFile.Filename = logC.Director + "/info.log"
			//}
			if logC.Format == "json" {
				infoHandler := slog.NewJSONHandler(infoFile, &opts)
				levelHandler.AddHandler(slog.LevelInfo, infoHandler)
			} else {
				infoHandler := slog.NewTextHandler(infoFile, &opts)
				levelHandler.AddHandler(slog.LevelInfo, infoHandler)
			}
		}
		if opts.Level.Level() <= slog.LevelWarn {

			// 始终创建文件写入器
			warnFile := &lumberjack.Logger{
				// 日志文件名，归档日志也会保存在对应目录下
				// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
				// <processname>-lumberjack.log
				//Filename: logC.Director + "/warn.log",

				// backup的日志是否使用本地时间戳，默认使用UTC时间
				LocalTime: true,
				// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
				MaxSize: logC.GetMaxSize(),
				// 旧日志保存的最大天数，默认保存所有旧日志文件
				MaxAge: logC.GetMaxAge(),
				// 旧日志保存的最大数量，默认保存所有旧日志文件
				MaxBackups: logC.GetMaxBackups(),
				// 对backup的日志是否进行压缩，默认不压缩
				Compress: true,
			}

			// if logC.WarnFile != "" {
			// 	warnFile.Filename = logC.Director + "/" + logC.WarnFile
			// } else {
			warnFile.Filename = logC.Director + "/warn.log"
			//}
			// warnHandler := slog.NewJSONHandler(warnFile, &opts)
			// levelHandler.AddHandler(slog.LevelWarn, warnHandler)
			if logC.Format == "json" {
				warnHandler := slog.NewJSONHandler(warnFile, &opts)
				levelHandler.AddHandler(slog.LevelWarn, warnHandler)
			} else {
				warnHandler := slog.NewTextHandler(warnFile, &opts)
				levelHandler.AddHandler(slog.LevelWarn, warnHandler)
			}
		}
		if opts.Level.Level() <= slog.LevelError {
			// 始终创建文件写入器
			errorFile := &lumberjack.Logger{
				// 日志文件名，归档日志也会保存在对应目录下
				// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
				// <processname>-lumberjack.log
				//Filename: logC.Director + "/error.log",

				// backup的日志是否使用本地时间戳，默认使用UTC时间
				LocalTime: true,
				// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
				MaxSize: logC.GetMaxSize(),
				// 旧日志保存的最大天数，默认保存所有旧日志文件
				MaxAge: logC.GetMaxAge(),
				// 旧日志保存的最大数量，默认保存所有旧日志文件
				MaxBackups: logC.GetMaxBackups(),
				// 对backup的日志是否进行压缩，默认不压缩
				Compress: true,
			}

			// if logC.ErrorFile != "" {
			// 	errorFile.Filename = logC.Director + "/" + logC.ErrorFile
			// } else {
			errorFile.Filename = logC.Director + "/error.log"
			//}
			if logC.Format == "json" {
				errorHandler := slog.NewJSONHandler(errorFile, &opts)
				levelHandler.AddHandler(slog.LevelError, errorHandler)
			} else {
				errorHandler := slog.NewTextHandler(errorFile, &opts)
				levelHandler.AddHandler(slog.LevelError, errorHandler)
			}
		}
	}

	if logC.LogInConsole {
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
	clog := slog.New(levelHandler)
	// 设置默认 logger
	slog.SetDefault(clog)
	return clog
}

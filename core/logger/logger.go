package logger

import (
	"io"
	"os"
	"time"

	"github.com/baowk/dilu-core/config"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

// Log 全局 logger 实例
var Log zerolog.Logger

// InitLogger 初始化 zerolog，返回配置好的 Logger
func InitLogger(logC config.LogCfg) zerolog.Logger {
	// 全局日志级别
	switch logC.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	filename := logC.Director + "/dilu.log"
	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		LocalTime:  true,
		MaxSize:    logC.GetMaxSize(),
		MaxAge:     logC.GetMaxAge(),
		MaxBackups: logC.GetMaxBackups(),
		Compress:   true,
	}

	var writers []io.Writer

	if logC.Format == "json" {
		writers = append(writers, fileWriter)
	} else {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        fileWriter,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
	}

	if logC.LogInConsole {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	var out io.Writer
	if len(writers) == 1 {
		out = writers[0]
	} else {
		out = zerolog.MultiLevelWriter(writers...)
	}

	Log = zerolog.New(out).With().Timestamp().Logger()
	return Log
}

package core

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"time"

// 	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// var Zap = new(_zap)

// type _zap struct{}

// // GetEncoder 获取 zapcore.Encoder
// func (z *_zap) GetEncoder() zapcore.Encoder {
// 	if Cfg.Logger.Format == "json" {
// 		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
// 	}
// 	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
// }

// // GetEncoderConfig 获取zapcore.EncoderConfig
// func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
// 	return zapcore.EncoderConfig{
// 		MessageKey:     "message",
// 		LevelKey:       "level",
// 		TimeKey:        "time",
// 		NameKey:        "logger",
// 		CallerKey:      "caller",
// 		StacktraceKey:  Cfg.Logger.StacktraceKey,
// 		LineEnding:     zapcore.DefaultLineEnding,
// 		EncodeLevel:    Cfg.Logger.ZapEncodeLevel(),
// 		EncodeTime:     z.CustomTimeEncoder,
// 		EncodeDuration: zapcore.SecondsDurationEncoder,
// 		EncodeCaller:   zapcore.FullCallerEncoder,
// 	}
// }

// // GetEncoderCore 获取Encoder的 zapcore.Core
// func (z *_zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
// 	filePath := path.Join(Cfg.Logger.Director, "%Y-%m-%d", l.String()+".log")

// 	w, err := GetWriter(filePath)
// 	if err != nil {
// 		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
// 		return nil
// 	}
// 	return zapcore.NewCore(z.GetEncoder(), w, level)
// }

// // 日志文件切割
// func GetWriter(filename string) (zapcore.WriteSyncer, error) {
// 	//保存日志30天，每1分钟分割一次日志
// 	hook, err := rotatelogs.New(
// 		filename,
// 		rotatelogs.WithClock(rotatelogs.Local),
// 		rotatelogs.WithMaxAge(24*time.Hour*time.Duration(Cfg.Logger.GetMaxAge())),
// 		rotatelogs.WithRotationTime(time.Hour*24),
// 	)
// 	if Cfg.Logger.LogInConsole {
// 		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), err
// 	}
// 	return zapcore.AddSync(hook), err
// }

// // CustomTimeEncoder 自定义日志输出时间格式
// func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
// 	if Cfg.Logger.Prefix == "" {
// 		encoder.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
// 	} else {
// 		encoder.AppendString(Cfg.Logger.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
// 	}
// }

// // GetZapCores 根据配置文件的Level获取 []zapcore.Core
// func (z *_zap) GetZapCores() []zapcore.Core {
// 	cores := make([]zapcore.Core, 0, 7)
// 	for level := Cfg.Logger.TransportLevel(); level <= zapcore.FatalLevel; level++ {
// 		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
// 	}
// 	return cores
// }

// // GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
// // Author [SliverHorn](https://github.com/SliverHorn)
// func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
// 	switch level {
// 	case zapcore.DebugLevel:
// 		return func(level zapcore.Level) bool { // 调试级别
// 			return level == zap.DebugLevel
// 		}
// 	case zapcore.InfoLevel:
// 		return func(level zapcore.Level) bool { // 日志级别
// 			return level == zap.InfoLevel
// 		}
// 	case zapcore.WarnLevel:
// 		return func(level zapcore.Level) bool { // 警告级别
// 			return level == zap.WarnLevel
// 		}
// 	case zapcore.ErrorLevel:
// 		return func(level zapcore.Level) bool { // 错误级别
// 			return level == zap.ErrorLevel
// 		}
// 	case zapcore.DPanicLevel:
// 		return func(level zapcore.Level) bool { // dpanic级别
// 			return level == zap.DPanicLevel
// 		}
// 	case zapcore.PanicLevel:
// 		return func(level zapcore.Level) bool { // panic级别
// 			return level == zap.PanicLevel
// 		}
// 	case zapcore.FatalLevel:
// 		return func(level zapcore.Level) bool { // 终止级别
// 			return level == zap.FatalLevel
// 		}
// 	default:
// 		return func(level zapcore.Level) bool { // 调试级别
// 			return level == zap.DebugLevel
// 		}
// 	}
// }

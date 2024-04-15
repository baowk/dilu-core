package config

// import (
// 	"strings"

// 	"go.uber.org/zap/zapcore"
// )

type LogCfg struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format       string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director     string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	MaxAge       int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间 天
	MaxSize      int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`                   // 日志文件大小
	MaxBackups   int    `mapstructure:"max-backups" json:"max-backups" yaml:"max-backups"`          // 日志备份天数
	ShowLine     bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
	//StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	EncodeLevel string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"` // 编码级
}

func (z *LogCfg) GetMaxAge() int {
	if z.MaxAge > 0 {
		return z.MaxAge
	}
	return 7
}

func (z *LogCfg) GetMaxSize() int {
	if z.MaxSize > 0 {
		return z.MaxSize
	}
	return 100
}

func (z *LogCfg) GetMaxBackups() int {
	if z.MaxBackups > 0 {
		return z.MaxBackups
	}
	return 7
}

// // ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
// func (z *LogCfg) ZapEncodeLevel() zapcore.LevelEncoder {
// 	switch {
// 	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
// 		return zapcore.LowercaseLevelEncoder
// 	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
// 		return zapcore.LowercaseColorLevelEncoder
// 	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
// 		return zapcore.CapitalLevelEncoder
// 	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
// 		return zapcore.CapitalColorLevelEncoder
// 	default:
// 		return zapcore.LowercaseLevelEncoder
// 	}
// }

func (z *LogCfg) Color() bool {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return false
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return true
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return false
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return true
	default:
		return false
	}
}

// // TransportLevel 根据字符串转化为 zapcore.Level
// func (z *LogCfg) TransportLevel() zapcore.Level {
// 	z.Level = strings.ToLower(z.Level)
// 	switch z.Level {
// 	case "debug":
// 		return zapcore.DebugLevel
// 	case "info":
// 		return zapcore.InfoLevel
// 	case "warn":
// 		return zapcore.WarnLevel
// 	case "error":
// 		return zapcore.WarnLevel
// 	case "dpanic":
// 		return zapcore.DPanicLevel
// 	case "panic":
// 		return zapcore.PanicLevel
// 	case "fatal":
// 		return zapcore.FatalLevel
// 	default:
// 		return zapcore.DebugLevel
// 	}
// }

package config

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
	EncodeLevel  string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	OutputMode   string `mapstructure:"output-mode" json:"output-mode" yaml:"output-mode"`          // 日志输出模式 level 按等级输出 single 单个文件
}

func (z *LogCfg) GetMaxAge() int {
	if z.MaxAge < 1 {
		z.MaxAge = 7
	}
	return z.MaxAge
}

func (z *LogCfg) GetMaxSize() int {
	if z.MaxSize < 1 {
		z.MaxSize = 100
	}
	return z.MaxSize
}

func (z *LogCfg) GetMaxBackups() int {
	if z.MaxBackups < 1 {
		z.MaxBackups = 7
	}
	return z.MaxBackups
}

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

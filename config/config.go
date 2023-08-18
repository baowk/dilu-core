package config

type AppCfg struct {
	Server ServerCfg         `mapstructure:"server" json:"server" yaml:"server"`
	Logger LogCfg            `mapstructure:"logger" json:"logger" yaml:"logger"`
	JWT    JWT               `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	DBCfg  DBCfg             `mapstructure:"dbcfg" json:"dbcfg" yaml:"dbcfg"` // 数据库配置
	Cache  CacheCfg          `mapstructure:"cache" json:"cache" yaml:"cache"` // 缓存
	Cors   CORS              `mapstructure:"cors" json:"cors" yaml:"cors"`
	Extend map[string]string `mapstructure:"extend" json:"extend" yaml:"extend"`
}

func (e *AppCfg) GetExtend(key string) string {
	return e.Extend[key]
}

type ServerCfg struct {
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"read-timeout" json:"read-timeout" yaml:"read-timeout"`
	WriteTimeout int    `mapstructure:"write-timeout" json:"write-timeout" yaml:"write-timeout"`
	FSType       string `mapstructure:"fs-type" json:"fs-type" yaml:"fs-type"`
	I18n         bool   `mapstructure:"i18n" json:"i18n" yaml:"i18n"` //是否开启多语言
	Lang         string `mapstructure:"lang" json:"lang" yaml:"lang"` //默认语言
}

func (e *ServerCfg) GetLang() string {
	if e.Lang == "" {
		return "zh-CN"
	}
	return e.Lang
}

func (e *ServerCfg) GetHost() string {
	if e.Host == "" {
		return "0.0.0.0"
	}
	return e.Host
}

func (e *ServerCfg) GetPort() int {
	if e.Port < 1 {
		return 7788
	}
	return e.Port
}

func (e *ServerCfg) GetReadTimeout() int {
	if e.ReadTimeout < 1 {
		return 20
	}
	return e.ReadTimeout
}

func (e *ServerCfg) GetWriteTimeout() int {
	if e.WriteTimeout < 1 {
		return 20
	}
	return e.WriteTimeout
}

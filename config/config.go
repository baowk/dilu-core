package config

type AppCfg struct {
	Server  ServerCfg `mapstructure:"server" json:"server" yaml:"server"` //服务器配置
	Logger  LogCfg    `mapstructure:"logger" json:"logger" yaml:"logger"` //log配置
	JWT     JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`          //jwt配置
	DBCfg   DBCfg     `mapstructure:"dbcfg" json:"dbcfg" yaml:"dbcfg"`    // 数据库配置
	Cache   CacheCfg  `mapstructure:"cache" json:"cache" yaml:"cache"`    // 缓存
	Cors    CORS      `mapstructure:"cors" json:"cors" yaml:"cors"`       //cors配置
	Extends Extend    `mapstructure:"extend" json:"extend" yaml:"extend"` //扩展配置
	Gen     bool      `mapstructure:"gen" json:"gen" yaml:"gen"`          //是否可生成初始化模式
}

type Extend map[string]any

func (e *Extend) Get(key string) any {
	return (*e)[key]
}

func (e *Extend) GetString(key string) string {
	if strVal, ok := (*e)[key].(string); ok {
		return strVal
	}
	return ""
}

func (e *Extend) GetInt(key string) int {
	if strVal, ok := (*e)[key].(int); ok {
		return strVal
	}
	return 0
}

func (e *Extend) GetFloat(key string) float64 {
	if strVal, ok := (*e)[key].(float64); ok {
		return strVal
	}
	return 0
}

func (e *Extend) GetBool(key string) bool {
	if strVal, ok := (*e)[key].(bool); ok {
		return strVal
	}
	return false
}

type ServerCfg struct {
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`                            //模式
	Host         string `mapstructure:"host" json:"host" yaml:"host"`                            //启动id
	Name         string `mapstructure:"name" json:"name" yaml:"name"`                            //appname
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`                            //端口
	ReadTimeout  int    `mapstructure:"read-timeout" json:"read-timeout" yaml:"read-timeout"`    //读超时
	WriteTimeout int    `mapstructure:"write-timeout" json:"write-timeout" yaml:"write-timeout"` //写超时
	FSType       string `mapstructure:"fs-type" json:"fs-type" yaml:"fs-type"`                   //文件系统
	I18n         bool   `mapstructure:"i18n" json:"i18n" yaml:"i18n"`                            //是否开启多语言
	Lang         string `mapstructure:"lang" json:"lang" yaml:"lang"`                            //默认语言
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

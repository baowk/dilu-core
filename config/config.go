package config

type AppCfg struct {
	Server  ServerCfg `mapstructure:"server" json:"server" yaml:"server"`    //服务器配置
	Remote  RemoteCfg `mapstructure:"remote" json:"remote" yaml:"remote"`    //远程配置
	Logger  LogCfg    `mapstructure:"logger" json:"logger" yaml:"logger"`    //log配置
	JWT     JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`             //jwt配置
	DBCfg   DBCfg     `mapstructure:"dbcfg" json:"dbcfg" yaml:"dbcfg"`       // 数据库配置
	Cache   CacheCfg  `mapstructure:"cache" json:"cache" yaml:"cache"`       // 缓存
	Cors    CORS      `mapstructure:"cors" json:"cors" yaml:"cors"`          //cors配置
	Extends any       `mapstructure:"extend" json:"extend" yaml:"extend"`    //扩展配置
	Gen     GenCfg    `mapstructure:"gen" json:"gen" yaml:"gen"`             //是否可生成
	Mongodb Mongodb   `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"` //mongo配置
}

type ServerCfg struct {
	Name         string `mapstructure:"name" json:"name" yaml:"name"`                            //appname
	RemoteEnable bool   `mapstructure:"remote-enable" json:"remote-enable" yaml:"remote-enable"` //是否开启远程配置
	Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`                            //模式
	Host         string `mapstructure:"host" json:"host" yaml:"host"`                            //启动host
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

type RemoteCfg struct {
	Provider      string `mapstructure:"provider" json:"provider" yaml:"provider"`                   //提供方
	Endpoint      string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`                   //端点
	Path          string `mapstructure:"path" json:"path" yaml:"path"`                               //路径
	SecretKeyring string `mapstructure:"secret-keyring" json:"secret-keyring" yaml:"secret-keyring"` //安全
	ConfigType    string `mapstructure:"config-type" json:"config-type" yaml:"config-type"`          //配置类型
}

func (e *RemoteCfg) GetConfigType() string {
	if e.ConfigType == "" {
		return "yaml"
	}
	return e.ConfigType
}

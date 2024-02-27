package config

type AppCfg struct {
	Server     ServerCfg     `mapstructure:"server" json:"server" yaml:"server"`                //服务器配置
	Remote     RemoteCfg     `mapstructure:"remote" json:"remote" yaml:"remote"`                //远程配置
	Logger     LogCfg        `mapstructure:"logger" json:"logger" yaml:"logger"`                //log配置
	JWT        JWT           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`                         //jwt配置
	DBCfg      DBCfg         `mapstructure:"dbcfg" json:"dbcfg" yaml:"dbcfg"`                   // 数据库配置
	Cache      CacheCfg      `mapstructure:"cache" json:"cache" yaml:"cache"`                   // 缓存
	Cors       CORS          `mapstructure:"cors" json:"cors" yaml:"cors"`                      //cors配置
	Extends    any           `mapstructure:"extend" json:"extend" yaml:"extend"`                //扩展配置
	Gen        GenCfg        `mapstructure:"gen" json:"gen" yaml:"gen"`                         //是否可生成
	GrpcServer GrpcServerCfg `mapstructure:"grpc-server" json:"grpc-server" yaml:"grpc-server"` //grpc服务配置
	// RdConfig   rd.Config     `mapstructure:"rd-config" json:"rd-config" yaml:"rd-config"`       //注册中心配置
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
	CloseWait    int    `mapstructure:"close-wait" json:"close-wait" yaml:"close-wait"`          //服务关闭等待 秒
	// Registry     bool   `mapstructure:"registry" json:"registry" yaml:"registry"`                //是否开启注册中心
}

type GrpcServerCfg struct {
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"` //启用Grpc服务
	Name   string `mapstructure:"name" json:"name" yaml:"name"`       //服务名，不设置默认为ServerName+"_grpc"
	Host   string `mapstructure:"host" json:"host" yaml:"host"`       //启动host
	Port   int    `mapstructure:"port" json:"port" yaml:"port"`       //端口
	//	Registry bool   `mapstructure:"registry" json:"registry" yaml:"registry"` //是否开启注册中心
	// Mode         string `mapstructure:"mode" json:"mode" yaml:"mode"`                            //模式
	// ReadTimeout  int    `mapstructure:"read-timeout" json:"read-timeout" yaml:"read-timeout"`    //读超时
	// WriteTimeout int    `mapstructure:"write-timeout" json:"write-timeout" yaml:"write-timeout"` //写超时
}

func (e *GrpcServerCfg) GetHost() string {
	if e.Host == "" {
		return "0.0.0.0"
	}
	return e.Host
}

func (e *GrpcServerCfg) GetPort() int {
	if e.Port < 1 {
		return 7789
	}
	return e.Port
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

func (e *ServerCfg) GetCloseWait() int {
	if e.CloseWait < 1 {
		return 1
	}
	return e.CloseWait
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

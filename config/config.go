package config

type Config interface {
	GetServerCfg() *ServerCfg
	GetLogCfg() *LogCfg
	GetDBCfg() *DBCfg
	GetCacheCfg() *CacheCfg
}

type ServerCfg struct {
	Name         string `mapstructure:"name" json:"name" yaml:"name"`                            //appname
	Node         int64  `mapstructure:"node" json:"node" yaml:"node"`                            //节点编号
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
}

func (e *ServerCfg) GetLang() string {
	if e.Lang == "" {
		e.Lang = "zh-CN"
	}
	return e.Lang
}

func (e *ServerCfg) GetHost() string {
	if e.Host == "" {
		e.Host = "0.0.0.0"
	}
	return e.Host
}

func (e *ServerCfg) GetPort() int {
	if e.Port < 1 {
		e.Port = 7788
	}
	return e.Port
}

func (e *ServerCfg) GetCloseWait() int {
	if e.CloseWait < 1 {
		e.CloseWait = 1
	}
	return e.CloseWait
}

func (e *ServerCfg) GetReadTimeout() int {
	if e.ReadTimeout < 1 {
		e.ReadTimeout = 20
	}
	return e.ReadTimeout
}

func (e *ServerCfg) GetWriteTimeout() int {
	if e.WriteTimeout < 1 {
		e.WriteTimeout = 20
	}
	return e.WriteTimeout
}

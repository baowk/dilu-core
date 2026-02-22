package config

type RemoteCfg struct {
	Provider      string `mapstructure:"provider" json:"provider" yaml:"provider"`                   //提供方
	Endpoint      string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`                   //端点
	Path          string `mapstructure:"path" json:"path" yaml:"path"`                               //路径
	SecretKeyring string `mapstructure:"secret-keyring" json:"secret-keyring" yaml:"secret-keyring"` //安全
	Token         string `mapstructure:"token" json:"token" yaml:"token"`                            //token
	ConfigType    string `mapstructure:"config-type" json:"config-type" yaml:"config-type"`          //配置类型
}

func (e *RemoteCfg) GetConfigType() string {
	if e.ConfigType == "" {
		e.ConfigType = "yaml"
	}
	return e.ConfigType
}

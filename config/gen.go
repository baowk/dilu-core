package config

type GenCfg struct {
	Enable    bool   `mapstructure:"enable" json:"enable" yaml:"ebable"`             // 开启生成
	FrontPath string `mapstructure:"front-path" json:"front-path" yaml:"front-path"` // 前端路径
}

package config

type Mongodb struct {
	URL     string `mapstructure:"url" json:"url" yaml:"url"`             // 链接地址
	Open    bool   `mapstructure:"open" json:"open" yaml:"open"`          // 开启状态
	Timeout int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"` // 超时时间
}

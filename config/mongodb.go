package config

type Mongodb struct {
	URL      string `mapstructure:"url" json:"url" yaml:"url"`                // 链接地址
	Database string `mapstructure:"database" json:"database" yaml:"database"` // 数据路
	Timeout  int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`    // 超时时间
}

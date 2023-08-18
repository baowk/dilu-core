package config

type JWT struct {
	SignKey string `mapstructure:"sign-key" json:"sign-key" yaml:"sign-key"` // jwt签名
	Expires int    `mapstructure:"expires" json:"expires" yaml:"expires"`    // 有效时长 分钟
	Refresh int    `mapstructure:"refresh" json:"refresh" yaml:"refresh"`    // 刷新时长
	Issuer  string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`       // 签发人
	Subject string `mapstructure:"subject" json:"subject" yaml:"subject"`    // 签发主体
}

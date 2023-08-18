package config

type CacheCfg struct {
	Type     string `mapstructure:"type" json:"type" yaml:"type"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

func (c *CacheCfg) GetType() string {
	if c.Type == "" {
		return "redis"
	}
	return c.Type
}

package config

type CacheCfg struct {
	Type       string `mapstructure:"type" json:"type" yaml:"type"`
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
	DB         int    `mapstructure:"db" json:"db" yaml:"db"`
	Prefix     string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	MasterName string `mapstructure:"master-name" json:"master-name" yaml:"master-name"`
}

func (c *CacheCfg) GetType() string {
	if c.Type == "" {
		c.Type = "memory"
	}
	return c.Type
}

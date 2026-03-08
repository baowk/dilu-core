package config

type GenCfg struct {
	Enable       bool              `mapstructure:"enable" json:"enable" yaml:"enable"`                      // 开启生成
	GenFront     bool              `mapstructure:"gen-front" json:"gen-front" yaml:"gen-front"`             // 是否生成前端代码
	FrontPath    string            `mapstructure:"front-path" json:"front-path" yaml:"front-path"`          // 前端路径
	TemplatePath string            `mapstructure:"template-path" json:"template-path" yaml:"template-path"` // 模板根目录
	ModuleMap    map[string]string `mapstructure:"module-map" json:"module-map" yaml:"module-map"`          // dbName -> moduleName
}

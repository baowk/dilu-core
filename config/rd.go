package config

import "time"

// RdCfg 注册发现配置（可选）
type RdCfg struct {
	Enable    bool          `mapstructure:"enable" json:"enable" yaml:"enable"`          // 是否启用
	Driver    string        `mapstructure:"driver" json:"driver" yaml:"driver"`          // etcd 或 consul
	Endpoints []string      `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints"` // 注册中心地址
	Scheme    string        `mapstructure:"scheme" json:"scheme" yaml:"scheme"`          // http/https
	Timeout   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`       // 超时
	Tags      []string      `mapstructure:"tags" json:"tags" yaml:"tags"`                // 服务标签
	Namespace string        `mapstructure:"namespace" json:"namespace" yaml:"namespace"` // 命名空间
}

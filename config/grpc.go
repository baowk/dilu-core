package config

type GrpcServerCfg struct {
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"` //启用Grpc服务
	Name   string `mapstructure:"name" json:"name" yaml:"name"`       //服务名，不设置默认为ServerName+"_grpc"
	Host   string `mapstructure:"host" json:"host" yaml:"host"`       //启动host
	Port   int    `mapstructure:"port" json:"port" yaml:"port"`       //端口
}

func (e *GrpcServerCfg) GetHost() string {
	if e.Host == "" {
		e.Host = "0.0.0.0"
	}
	return e.Host
}

func (e *GrpcServerCfg) GetPort() int {
	if e.Port < 1 {
		e.Port = 7789
	}
	return e.Port
}

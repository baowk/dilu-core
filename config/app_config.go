package config

import (
	"sync"
)

// AppConfig 实现Config接口的具体配置结构
type AppConfig struct {
	Server      ServerCfg     `mapstructure:"server" json:"server" yaml:"server"`
	Remote      RemoteCfg     `mapstructure:"remote" json:"remote" yaml:"remote"`
	Logger      LogCfg        `mapstructure:"logger" json:"logger" yaml:"logger"`
	JWT         JWT           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	DBCfg       DBCfg         `mapstructure:"dbcfg" json:"dbcfg" yaml:"dbcfg"`
	Cache       CacheCfg      `mapstructure:"cache" json:"cache" yaml:"cache"`
	Cors        CORS          `mapstructure:"cors" json:"cors" yaml:"cors"`
	Gen         GenCfg        `mapstructure:"gen" json:"gen" yaml:"gen"`
	GrpcServer  GrpcServerCfg `mapstructure:"grpc-server" json:"grpc-server" yaml:"grpc-server"`
	AccessLimit AccessLimit   `mapstructure:"access-limit" json:"access-limit" yaml:"access-limit"`
	Rd          RdCfg         `mapstructure:"rd" json:"rd" yaml:"rd"`
}

// 确保AppConfig实现Config接口
var _ Config = (*AppConfig)(nil)

func (a *AppConfig) GetServerCfg() *ServerCfg {
	return &a.Server
}

func (a *AppConfig) GetLogCfg() *LogCfg {
	return &a.Logger
}

func (a *AppConfig) GetDBCfg() *DBCfg {
	return &a.DBCfg
}

func (a *AppConfig) GetCacheCfg() *CacheCfg {
	return &a.Cache
}

// ConfigManager 配置管理器，支持热加载和线程安全
type ConfigManager struct {
	config Config
	mu     sync.RWMutex
}

func NewConfigManager(initialConfig Config) *ConfigManager {
	return &ConfigManager{
		config: initialConfig,
	}
}

func (cm *ConfigManager) GetConfig() Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

func (cm *ConfigManager) SetConfig(newConfig Config) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = newConfig
}

func (cm *ConfigManager) GetServerCfg() *ServerCfg {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.GetServerCfg()
}

func (cm *ConfigManager) GetLogCfg() *LogCfg {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.GetLogCfg()
}

func (cm *ConfigManager) GetDBCfg() *DBCfg {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.GetDBCfg()
}

func (cm *ConfigManager) GetCacheCfg() *CacheCfg {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.GetCacheCfg()
}
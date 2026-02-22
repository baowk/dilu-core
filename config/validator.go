package config

import (
	"fmt"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	errors []ValidationError
}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		errors: make([]ValidationError, 0),
	}
}

// ValidationError 配置验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Validate 验证配置
func (cv *ConfigValidator) Validate(cfg Config) []ValidationError {
	cv.errors = make([]ValidationError, 0)
	
	// 验证服务器配置
	cv.validateServerConfig(cfg.GetServerCfg())
	
	// 验证日志配置
	cv.validateLogConfig(cfg.GetLogCfg())
	
	// 验证数据库配置
	cv.validateDBConfig(cfg.GetDBCfg())
	
	// 验证缓存配置
	cv.validateCacheConfig(cfg.GetCacheCfg())
	
	return cv.errors
}

// validateServerConfig 验证服务器配置
func (cv *ConfigValidator) validateServerConfig(server *ServerCfg) {
	if server == nil {
		cv.addError("server", "server configuration is required", nil)
		return
	}

	// 验证端口
	if server.Port < 1 || server.Port > 65535 {
		cv.addError("server.port", "port must be between 1 and 65535", server.Port)
	}

	// 验证主机
	if server.Host == "" {
		cv.addError("server.host", "host is required", server.Host)
	}

	// 验证名称
	if server.Name == "" {
		cv.addError("server.name", "name is required", server.Name)
	}

	// 验证模式
	validModes := map[string]bool{"dev": true, "test": true, "prod": true}
	if server.Mode != "" && !validModes[server.Mode] {
		cv.addError("server.mode", "mode must be one of: dev, test, prod", server.Mode)
	}

	// 验证超时设置
	if server.ReadTimeout < 0 {
		cv.addError("server.read_timeout", "read timeout cannot be negative", server.ReadTimeout)
	}
	if server.WriteTimeout < 0 {
		cv.addError("server.write_timeout", "write timeout cannot be negative", server.WriteTimeout)
	}
	if server.CloseWait < 0 {
		cv.addError("server.close_wait", "close wait cannot be negative", server.CloseWait)
	}
}

// validateLogConfig 验证日志配置
func (cv *ConfigValidator) validateLogConfig(log *LogCfg) {
	if log == nil {
		cv.addError("logger", "logger configuration is required", nil)
		return
	}

	// 验证日志级别
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if log.Level != "" && !validLevels[log.Level] {
		cv.addError("logger.level", "level must be one of: debug, info, warn, error", log.Level)
	}

	// 验证目录
	if log.Director == "" {
		cv.addError("logger.director", "director is required", log.Director)
	}

	// 验证文件大小限制
	if log.MaxSize < 0 {
		cv.addError("logger.max_size", "max size cannot be negative", log.MaxSize)
	}

	// 验证保留天数
	if log.MaxAge < 0 {
		cv.addError("logger.max_age", "max age cannot be negative", log.MaxAge)
	}

	// 验证备份数量
	if log.MaxBackups < 0 {
		cv.addError("logger.max_backups", "max backups cannot be negative", log.MaxBackups)
	}
}

// validateDBConfig 验证数据库配置
func (cv *ConfigValidator) validateDBConfig(db *DBCfg) {
	if db == nil {
		cv.addError("dbcfg", "database configuration is required", nil)
		return
	}

	// 验证主数据库配置
	if db.DSN == "" && len(db.DBS) == 0 {
		cv.addError("dbcfg", "either main database DSN or additional databases must be configured", nil)
	}

	// 如果有主数据库DSN，验证其配置
	if db.DSN != "" {
		cv.validateMainDBConfig(db)
	}

	// 验证额外数据库配置
	for key, dbc := range db.DBS {
		fieldPrefix := fmt.Sprintf("dbcfg.dbs.%s", key)
		cv.validateSingleDBConfig(fieldPrefix, dbc.DSN, dbc.Driver, dbc.MaxIdleConns, dbc.MaxOpenConns, dbc.MaxLifetime)
	}
}

// validateMainDBConfig 验证主数据库配置
func (cv *ConfigValidator) validateMainDBConfig(db *DBCfg) {
	// 验证驱动
	validDrivers := map[string]bool{"mysql": true, "postgres": true, "sqlite": true, "clickhouse": true}
	if db.Driver != "" && !validDrivers[db.Driver] {
		cv.addError("dbcfg.driver", "driver must be one of: mysql, postgres, sqlite, clickhouse", db.Driver)
	}

	// 验证连接池设置
	if db.MaxIdleConns < 0 {
		cv.addError("dbcfg.max_idle_conns", "max idle connections cannot be negative", db.MaxIdleConns)
	}
	if db.MaxOpenConns < 0 {
		cv.addError("dbcfg.max_open_conns", "max open connections cannot be negative", db.MaxOpenConns)
	}
	if db.MaxLifetime < 0 {
		cv.addError("dbcfg.max_lifetime", "max lifetime cannot be negative", db.MaxLifetime)
	}
}

// validateSingleDBConfig 验证单个数据库配置
func (cv *ConfigValidator) validateSingleDBConfig(prefix, dsn, driver string, maxIdle, maxOpen, maxLifetime int) {
	// 验证DSN
	if dsn == "" {
		cv.addError(prefix+".dsn", "DSN is required", dsn)
	}

	// 验证驱动
	validDrivers := map[string]bool{"mysql": true, "postgres": true, "sqlite": true, "clickhouse": true}
	if driver != "" && !validDrivers[driver] {
		cv.addError(prefix+".driver", "driver must be one of: mysql, postgres, sqlite, clickhouse", driver)
	}

	// 验证连接池设置
	if maxIdle < 0 {
		cv.addError(prefix+".max_idle_conns", "max idle connections cannot be negative", maxIdle)
	}
	if maxOpen < 0 {
		cv.addError(prefix+".max_open_conns", "max open connections cannot be negative", maxOpen)
	}
	if maxLifetime < 0 {
		cv.addError(prefix+".max_lifetime", "max lifetime cannot be negative", maxLifetime)
	}
}

// validateCacheConfig 验证缓存配置
func (cv *ConfigValidator) validateCacheConfig(cache *CacheCfg) {
	if cache == nil {
		cv.addError("cache", "cache configuration is required", nil)
		return
	}

	// 验证缓存类型
	validTypes := map[string]bool{"memory": true, "redis": true}
	if cache.Type != "" && !validTypes[cache.Type] {
		cv.addError("cache.type", "type must be one of: memory, redis", cache.Type)
	}

	// 如果是Redis类型，验证必要配置
	if cache.Type == "redis" {
		if cache.Addr == "" {
			cv.addError("cache.addr", "address is required for redis cache", cache.Addr)
		}
		if cache.DB < 0 {
			cv.addError("cache.db", "database number cannot be negative", cache.DB)
		}
	}
}

// addError 添加验证错误
func (cv *ConfigValidator) addError(field, message string, value interface{}) {
	cv.errors = append(cv.errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors 检查是否有验证错误
func (cv *ConfigValidator) HasErrors() bool {
	return len(cv.errors) > 0
}

// GetErrors 获取所有验证错误
func (cv *ConfigValidator) GetErrors() []ValidationError {
	return cv.errors
}

// ApplyDefaults 应用默认值
func ApplyDefaults(cfg Config) {
	server := cfg.GetServerCfg()
	if server != nil {
		applyServerDefaults(server)
	}

	log := cfg.GetLogCfg()
	if log != nil {
		applyLogDefaults(log)
	}

	db := cfg.GetDBCfg()
	if db != nil {
		applyDBDefaults(db)
	}

	cache := cfg.GetCacheCfg()
	if cache != nil {
		applyCacheDefaults(cache)
	}
}

// applyServerDefaults 应用服务器默认值
func applyServerDefaults(server *ServerCfg) {
	if server.Host == "" {
		server.Host = "0.0.0.0"
	}
	if server.Port == 0 {
		server.Port = 7788
	}
	if server.Mode == "" {
		server.Mode = "dev"
	}
	if server.ReadTimeout == 0 {
		server.ReadTimeout = 20
	}
	if server.WriteTimeout == 0 {
		server.WriteTimeout = 20
	}
	if server.CloseWait == 0 {
		server.CloseWait = 5
	}
	if server.Lang == "" {
		server.Lang = "zh-CN"
	}
}

// applyLogDefaults 应用日志默认值
func applyLogDefaults(log *LogCfg) {
	if log.Level == "" {
		log.Level = "info"
	}
	if log.Format == "" {
		log.Format = "json"
	}
	if log.Director == "" {
		log.Director = "logs"
	}
	if log.EncodeLevel == "" {
		log.EncodeLevel = "LowercaseLevelEncoder"
	}
	if log.MaxSize == 0 {
		log.MaxSize = 100
	}
	if log.MaxAge == 0 {
		log.MaxAge = 7
	}
	if log.MaxBackups == 0 {
		log.MaxBackups = 10
	}
	if log.OutputMode == "" {
		log.OutputMode = "level"
	}
}

// applyDBDefaults 应用数据库默认值
func applyDBDefaults(db *DBCfg) {
	if db.MaxIdleConns == 0 {
		db.MaxIdleConns = 10
	}
	if db.MaxOpenConns == 0 {
		db.MaxOpenConns = 100
	}
	if db.MaxLifetime == 0 {
		db.MaxLifetime = 30
	}
	if db.LogMode == "" {
		db.LogMode = "warn"
	}
	if db.SlowThreshold == 0 {
		db.SlowThreshold = 200
	}
	
	// 应用额外数据库的默认值
	for key := range db.DBS {
		dbc := db.DBS[key]
		if dbc.MaxIdleConns == 0 {
			dbc.MaxIdleConns = db.MaxIdleConns
		}
		if dbc.MaxOpenConns == 0 {
			dbc.MaxOpenConns = db.MaxOpenConns
		}
		if dbc.MaxLifetime == 0 {
			dbc.MaxLifetime = db.MaxLifetime
		}
		if dbc.LogMode == "" {
			dbc.LogMode = db.LogMode
		}
		if dbc.SlowThreshold == 0 {
			dbc.SlowThreshold = db.SlowThreshold
		}
		db.DBS[key] = dbc
	}
}

// applyCacheDefaults 应用缓存默认值
func applyCacheDefaults(cache *CacheCfg) {
	if cache.Type == "" {
		cache.Type = "memory"
	}
	if cache.DB == 0 {
		cache.DB = 0
	}
	if cache.Prefix == "" {
		cache.Prefix = "dilu:"
	}
}
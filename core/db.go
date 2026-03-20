package core

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/config"
	diluLogger "github.com/baowk/dilu-core/core/logger"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DBManager 数据库管理器
type DBManager struct {
	databases map[string]*gorm.DB
	logger    zerolog.Logger
}

// NewDBManager 创建数据库管理器
func NewDBManager(logger zerolog.Logger) *DBManager {
	return &DBManager{
		databases: make(map[string]*gorm.DB),
		logger:    logger,
	}
}

// dbInit 优化的数据库初始化方法
func (app *Application) dbInit() error {
	cfg := app.config
	dbManager := NewDBManager(diluLogger.Log)

	// 初始化日志写入器
	logWriter, err := app.createLogWriter(cfg)
	if err != nil {
		return fmt.Errorf("failed to create log writer: %w", err)
	}

	// 初始化主数据库
	if cfg.GetDBCfg().DSN != "" {
		logMode := config.GetLogMode(cfg.GetDBCfg().LogMode)
		db, err := dbManager.initDatabase(
			cfg.GetDBCfg().Driver, cfg.GetDBCfg().DSN, cfg.GetDBCfg().Prefix,
			consts.DB_DEF, logMode, cfg.GetDBCfg().SlowThreshold,
			cfg.GetDBCfg().MaxIdleConns, cfg.GetDBCfg().MaxOpenConns,
			cfg.GetDBCfg().MaxLifetime, cfg.GetDBCfg().Singular,
			cfg.GetLogCfg().Color(), cfg.GetDBCfg().IgnoreNotFound, logWriter,
		)
		if err != nil {
			return fmt.Errorf("failed to initialize main database: %w", err)
		}
		app.databases[consts.DB_DEF] = db
	}

	// 初始化额外数据库
	for key, dbc := range cfg.GetDBCfg().DBS {
		if !dbc.Disable {
			if err := app.initAdditionalDB(dbManager, key, dbc, cfg, logWriter); err != nil {
				diluLogger.Log.Error().Str("key", key).Err(err).Msg("Failed to initialize additional database")
				continue
			}
		}
	}

	return nil
}

// createLogWriter 创建日志写入器
func (app *Application) createLogWriter(cfg config.Config) (io.Writer, error) {
	logCfg := cfg.GetLogCfg()

	fileWriter := &lumberjack.Logger{
		Filename:   logCfg.Director + "/sql.log",
		LocalTime:  true,
		MaxSize:    logCfg.GetMaxSize(),
		MaxAge:     logCfg.GetMaxAge(),
		MaxBackups: logCfg.GetMaxBackups(),
		Compress:   true,
	}

	if logCfg.OutputMode == "single" {
		fileWriter.Filename = logCfg.Director + "/dilu.log"
	}

	if logCfg.LogInConsole {
		return io.MultiWriter(fileWriter, os.Stdout), nil
	}
	return fileWriter, nil
}

// initAdditionalDB 初始化额外数据库
func (app *Application) initAdditionalDB(
	dbManager *DBManager, key string, dbc config.DB,
	cfg config.Config, logWriter io.Writer,
) error {
	// 获取配置参数，使用默认值填充
	logMode := app.getLogMode(dbc.LogMode, cfg.GetDBCfg().LogMode)
	prefix := app.getPrefix(dbc.Prefix, cfg.GetDBCfg().Prefix)
	slowThreshold := app.getSlowThreshold(dbc.SlowThreshold, cfg.GetDBCfg().SlowThreshold)
	maxIdleConns := app.getMaxIdleConns(dbc.MaxIdleConns, cfg.GetDBCfg().MaxIdleConns)
	maxOpenConns := app.getMaxOpenConns(dbc.MaxOpenConns, cfg.GetDBCfg().MaxOpenConns)
	maxLifetime := app.getMaxLifetime(dbc.MaxLifetime, cfg.GetDBCfg().MaxLifetime)
	driver := app.getDriver(dbc.Driver, cfg.GetDBCfg().Driver)
	ignoreNotFound := dbc.IgnoreNotFound || cfg.GetDBCfg().IgnoreNotFound

	// 初始化数据库连接
	db, err := dbManager.initDatabase(
		driver, dbc.DSN, prefix, key, logMode, slowThreshold,
		maxIdleConns, maxOpenConns, maxLifetime,
		cfg.GetDBCfg().Singular, cfg.GetLogCfg().Color(),
		ignoreNotFound, logWriter,
	)
	if err != nil {
		return err
	}

	app.mu.Lock()
	app.databases[key] = db
	app.mu.Unlock()
	return nil
}

// 辅助方法获取配置参数
func (app *Application) getLogMode(dbMode, defaultMode string) logger.LogLevel {
	if dbMode != "" {
		return config.GetLogMode(dbMode)
	}
	return config.GetLogMode(defaultMode)
}

func (app *Application) getPrefix(dbPrefix, defaultPrefix string) string {
	if dbPrefix != "" {
		return dbPrefix
	}
	return defaultPrefix
}

func (app *Application) getSlowThreshold(dbSlow, defaultSlow int) int {
	if dbSlow > 0 {
		return dbSlow
	}
	return defaultSlow
}

func (app *Application) getMaxIdleConns(dbMaxIdle, defaultMaxIdle int) int {
	if dbMaxIdle > 0 {
		return dbMaxIdle
	}
	return defaultMaxIdle
}

func (app *Application) getMaxOpenConns(dbMaxOpen, defaultMaxOpen int) int {
	if dbMaxOpen > 0 {
		return dbMaxOpen
	}
	return defaultMaxOpen
}

func (app *Application) getMaxLifetime(dbMaxLifetime, defaultMaxLifetime int) int {
	if dbMaxLifetime > 0 {
		return dbMaxLifetime
	}
	return defaultMaxLifetime
}

func (app *Application) getDriver(dbDriver, defaultDriver string) string {
	if dbDriver != "" {
		return dbDriver
	}
	return defaultDriver
}

// initDatabase 核心数据库初始化逻辑
func (dm *DBManager) initDatabase(
	driver, dsn, prefix, key string, logMode logger.LogLevel,
	slowThreshold, maxIdleConns, maxOpenConns, maxLifetime int,
	singular, color bool, ignoreNotFound bool, logWriter io.Writer,
) (*gorm.DB, error) {
	// 选择数据库驱动
	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "clickhouse":
		dialector = clickhouse.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
	}

	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %s: %w", key, err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database %s: %w", key, err)
	}

	dm.logger.Info().Str("key", key).Str("driver", driver).Msg("Database connected successfully")
	dm.databases[key] = db
	return db, nil
}

// GetDB 获取数据库连接
func (dm *DBManager) GetDB(key string) (*gorm.DB, error) {
	db, exists := dm.databases[key]
	if !exists {
		return nil, fmt.Errorf("database %s not found", key)
	}
	return db, nil
}

// Close 关闭所有数据库连接
func (dm *DBManager) Close() error {
	var lastErr error
	for key, db := range dm.databases {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				lastErr = fmt.Errorf("failed to close database %s: %w", key, err)
			}
		}
	}
	return lastErr
}

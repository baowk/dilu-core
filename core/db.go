package core

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/baowk/dilu-core/common/consts"
	"github.com/baowk/dilu-core/config"
	"github.com/natefinch/lumberjack"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func dbInit() {

	// 始终创建文件写入器
	fileWriter := &lumberjack.Logger{
		// 日志文件名，归档日志也会保存在对应目录下
		// 若该值为空，则日志会保存到os.TempDir()目录下，日志文件名为
		// <processname>-lumberjack.log
		Filename: _cfg.GetLogCfg().Director + "/sql.log",

		// backup的日志是否使用本地时间戳，默认使用UTC时间
		LocalTime: true,
		// 日志大小到达MaxSize(MB)就开始backup，默认值是100.
		MaxSize: _cfg.GetLogCfg().GetMaxSize(),
		// 旧日志保存的最大天数，默认保存所有旧日志文件
		MaxAge: _cfg.GetLogCfg().GetMaxAge(),
		// 旧日志保存的最大数量，默认保存所有旧日志文件
		MaxBackups: _cfg.GetLogCfg().GetMaxBackups(),
		// 对backup的日志是否进行压缩，默认不压缩
		Compress: true,
	}

	if _cfg.GetLogCfg().OutputMode == "single" {
		fileWriter.Filename = _cfg.GetLogCfg().Director + "/dilu.log"
	}

	var logWrite io.Writer

	if _cfg.GetLogCfg().LogInConsole {
		// 同时输出到文件和控制台
		logWrite = io.MultiWriter(fileWriter, os.Stdout)
	} else {
		// 仅输出到文件
		logWrite = fileWriter
	}

	if _cfg.GetDBCfg().DSN != "" {
		logMode := config.GetLogMode(_cfg.GetDBCfg().LogMode)
		initDb(_cfg.GetDBCfg().Driver, _cfg.GetDBCfg().DSN, _cfg.GetDBCfg().Prefix, consts.DB_DEF, logMode, _cfg.GetDBCfg().SlowThreshold,
			_cfg.GetDBCfg().MaxIdleConns, _cfg.GetDBCfg().MaxOpenConns, _cfg.GetDBCfg().MaxLifetime, _cfg.GetDBCfg().Singular, _cfg.GetLogCfg().Color(), _cfg.GetDBCfg().IgnoreNotFound, logWrite)
	}
	for key, dbc := range _cfg.GetDBCfg().DBS {
		if !dbc.Disable {
			var logMode logger.LogLevel
			if dbc.LogMode != "" {
				logMode = config.GetLogMode(dbc.LogMode)
			} else {
				logMode = config.GetLogMode(_cfg.GetDBCfg().LogMode)
			}
			prefix := dbc.Prefix
			if prefix == "" && _cfg.GetDBCfg().Prefix != "" {
				prefix = _cfg.GetDBCfg().Prefix
			}
			slow := dbc.SlowThreshold
			if slow < 1 && _cfg.GetDBCfg().SlowThreshold > 0 {
				slow = _cfg.GetDBCfg().SlowThreshold
			}
			singular := _cfg.GetDBCfg().Singular
			maxIdle := dbc.MaxIdleConns
			if maxIdle < 1 {
				maxIdle = _cfg.GetDBCfg().GetMaxIdleConns()
			}

			maxOpen := dbc.MaxOpenConns
			if maxOpen < 1 {
				maxOpen = _cfg.GetDBCfg().GetMaxOpenConns()
			}

			maxLifetime := dbc.MaxLifetime
			if maxLifetime < 1 {
				maxLifetime = _cfg.GetDBCfg().GetMaxLifetime()
			}
			driver := dbc.Driver
			if driver == "" && _cfg.GetDBCfg().Driver != "" {
				driver = _cfg.GetDBCfg().Driver
			}
			ignoreNotFound := dbc.IgnoreNotFound
			if !ignoreNotFound && _cfg.GetDBCfg().IgnoreNotFound {
				ignoreNotFound = _cfg.GetDBCfg().IgnoreNotFound
			}
			initDb(driver, dbc.DSN, prefix, key, logMode, slow, maxIdle, maxOpen, maxLifetime, singular, _cfg.GetLogCfg().Color(), ignoreNotFound, logWrite)
		}
	}

}

func initDb(driver, dns, prefix, key string, logMode logger.LogLevel, slow, maxIdle, maxOpen, maxLifetime int, singular, color, ignoreNotFound bool, logWrite io.Writer) {
	var db *gorm.DB
	var err error
	switch driver {
	case Mysql.String():
		db, err = gorm.Open(mysql.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound, logWrite))
	case Pgsql.String():
		db, err = gorm.Open(postgres.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound, logWrite))
	case Sqlite.String():
		db, err = gorm.Open(sqlite.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound, logWrite))
	// case Mssql.String():
	// 	db, err = gorm.Open(sqlserver.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound, logWrite))
	// case "oracle":
	// 	db, err = gorm.Open(oracle.Open(dbc.DSN), &gorm.Config{})
	case ClickHouse.String():
		db, err = gorm.Open(clickhouse.Open(dns), GetGromLogCfg(logMode, prefix, slow, singular, color, ignoreNotFound, logWrite))
	default:
		err = errors.New("db err")
	}
	if err != nil {
		slog.Error("connect db err ", "dns", dns, "key", key, "err", err)
		panic(err)
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		slog.Error("connect db err ", "dns", dns, "key", key, "err", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(maxLifetime))
	SetDb(key, db)
	dbInitFlag = true
}

func GetGromLogCfg(logMode logger.LogLevel, prefix string, slowThreshold int, singular, color, ignoreNotFound bool, logW io.Writer) *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		//DisableForeignKeyConstraintWhenMigrating: true,
	}

	//filePath := path.Join(_cfg.GetLogCfg().Director, "%Y-%m-%d", "sql.log")
	//w, _ := GetWriter(filePath)
	slow := time.Duration(slowThreshold) * time.Millisecond
	_default := logger.New(log.New(logW, prefix, log.LstdFlags), logger.Config{
		SlowThreshold:             slow,
		Colorful:                  color,
		IgnoreRecordNotFoundError: ignoreNotFound,
	})

	config.Logger = _default.LogMode(logMode)

	return config
}

func SetDb(key string, db *gorm.DB) {
	// lock.Lock()
	// defer lock.Unlock()
	dbs[key] = db
}

// GetDb 获取所有map里的db数据
func Dbs() map[string]*gorm.DB {
	// lock.RLock()
	// defer lock.RUnlock()
	return dbs
}

func Db(name string) *gorm.DB {
	// lock.RLock()
	// defer lock.RUnlock()
	if dbInitFlag {
		if len(dbs) == 1 {
			return dbs[consts.DB_DEF]
		}
		if db, ok := dbs[name]; !ok || db == nil {
			slog.Error("db init err", "err", errors.New(name))
			panic("db not init")
		} else {
			return db
		}
	} else {
		return nil
	}
}

// 获取默认的（master）db
func DB() *gorm.DB {
	return Db(consts.DB_DEF)
}

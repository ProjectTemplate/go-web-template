package mysql

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

// dbMap 存储初始化后的数据库实例
var dbMap map[string]*gorm.DB

// GetDB 根据别名获取数据库连接 [gorm.DB]，别名必须存在，否则会panic
func GetDB(ctx context.Context, name string) *gorm.DB {
	db := dbMap[name]
	if db == nil {
		logger.Error(ctx, "GetDB error, db is nil", zap.String("name", name))
		panic("GetDB error, db is nil")
	}
	return db
}

// Init 根据配置信息初始化数据库连接，如果初始化失败会 panic，所以要确保配置信息都是正确的
//
// 配置信息说明
// 前缀 mysql 为固定前缀，后面的 test 为数据库连接的别名，GetDB方法使用别名获取对应的连接
// [mysql.test]
//
//	dsn = [#数据源的连接信息，支持多个，第一个库为主库，其余的为只读库
//	    "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local&parseTime=True",
//	    "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local&parseTime=True"
//	]
//	max_open_connections = 50  最大打开的连接数
//	max_idle_connections = 25  最大空闲连接数
//	max_life_time = "1h"       连接的最大存活时间
//	max_idle_time = "10m"      连接的最大空闲时间
//	show_log = false           是否展示MySQL日志
//	slow_threshold = "1ms"     慢查询阈值
//
//nolint:gocyclo
func Init(ctx context.Context, dbConfigs map[string]config.MySQL) {
	logger.Info(ctx, "init MySQL, config info: ", zap.Any("config", dbConfigs))
	dbMap = make(map[string]*gorm.DB)

	for dnName, dbConfig := range dbConfigs {
		logger.Info(ctx, "init single mysql connection, config info: ", zap.String("dnName", dnName), zap.Any("config", dbConfig))
		if len(dbConfig.DSN) < 1 {
			panic("init MySQL config error, dsn is empty. db name: " + dnName)
		}

		customGormLogger := NewGormLogger(dbConfig.SlowThreshold)
		if !dbConfig.ShowLog {
			customGormLogger.LogMode(gormLogger.Silent)
		} else {
			customGormLogger.LogMode(gormLogger.Info)
		}

		//主库，第一个配置为主库
		db, err := gorm.Open(mysql.Open(dbConfig.DSN[0]), &gorm.Config{
			Logger: customGormLogger,
		})

		if err != nil {
			logger.Error(ctx, "init single mysql connection, error", zap.String("dnName", dnName), zap.Error(err))
			panic("init MySQL, init single mysql connection. db name: " + dnName + ", error: " + err.Error())
		}

		//从库，除了第一个库，其余的库为从库
		var replicas = make([]gorm.Dialector, 0, len(dbConfig.DSN)-1)
		if len(dbConfig.DSN) > 1 {
			for i := range dbConfig.DSN[1:] {
				replicas = append(replicas, mysql.Open(dbConfig.DSN[i+1]))
			}
		}

		plugin := dbresolver.Register(dbresolver.Config{
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		})

		if dbConfig.MaxOpenConnections > 0 {
			plugin.SetMaxOpenConns(dbConfig.MaxOpenConnections)
		}

		if dbConfig.MaxIdleConnections > 0 {
			plugin.SetMaxIdleConns(dbConfig.MaxIdleConnections)
		}

		if dbConfig.MaxLifeTime > 0 {
			plugin.SetConnMaxLifetime(dbConfig.MaxLifeTime)
		}

		if dbConfig.MaxIdleTime > 0 {
			plugin.SetConnMaxIdleTime(dbConfig.MaxIdleTime)
		}

		err = db.Use(plugin)
		if err != nil {
			logger.Error(ctx, "init single mysql connection, error", zap.String("dnName", dnName), zap.Error(err))
			panic("init MySQL, init single mysql connection. db name: " + dnName + ", error: " + err.Error())
		}

		dbMap[dnName] = db

		logger.Info(ctx, "init single mysql connection success", zap.String("dnName", dnName))
	}
}

package mysql

import (
	"context"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// dbMap 存储初始化后的数据库实例
var dbMap map[string]*gorm.DB

// GetDB 根据名字获取数据库连接（db.[name]），名字必须存在，否则会panic
//
// 如果获取失败，会打印错误日志，并且panic，通过 panic 提示错误配置
func GetDB(name string) *gorm.DB {
	db := dbMap[name]
	if db == nil {
		logger.Error(context.Background(), "GetDB error, db is nil", zap.String("name", name))
		panic("GetDB error, db is nil")
	}
	return db
}

// Init 初始化数据库连接
//
// 在配置文件中 db.test.dsn 列表里面的第一个为主库，其余的为从库
// 完整配置文件参考 ./data/config.toml 文件
func Init(ctx context.Context, dbConfigs map[string]config.DB) {
	logger.Info(ctx, "init MySQL, config info: ", zap.Any("config", dbConfigs))
	dbMap = make(map[string]*gorm.DB)

	for dnName, dbConfig := range dbConfigs {
		logger.Info(ctx, "init single mysql connection, config info: ", zap.String("dnName", dnName), zap.Any("config", dbConfig))
		if len(dbConfig.DSN) < 1 {
			panic("init MySQL config error, dsn is empty. db name: " + dnName)
		}

		//主库，第一个配置为主库
		customGormLogger := NewGormLogger(dbConfig.SlowThreshold)
		if !dbConfig.ShowLog {
			customGormLogger.LogMode(gormLogger.Silent)
		} else {
			customGormLogger.LogMode(gormLogger.Info)
		}

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

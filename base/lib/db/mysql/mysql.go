package mysql

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

func Init(ctx context.Context, dbConfigs map[string]config.DB) (map[string]gorm.DB, error) {
	logger.Info(ctx, "init mysql, config info: ", zap.Any("config", dbConfigs))

	result := make(map[string]gorm.DB, len(dbConfigs))

	for dnName, dbConfig := range dbConfigs {
		logger.Info(ctx, "init single mysql connection, config info: ", zap.String("dnName", dnName), zap.Any("config", dbConfig))
		if len(dbConfig.DSN) < 1 {
			return result, errors.New("mysql config error, dsn is empty. db name: " + dnName)
		}

		dbLoggerLevel := gormLogger.Silent
		if dbConfig.IsLogger {
			dbLoggerLevel = gormLogger.Info
		}

		//主库，第一个配置为主库
		customGormLogger := NewGormLogger(dbConfig.SlowThreshold)
		customGormLogger.LogMode(dbLoggerLevel)

		db, err := gorm.Open(mysql.Open(dbConfig.DSN[0]), &gorm.Config{
			Logger: customGormLogger,
		})

		if err != nil {
			logger.Error(ctx, "init single mysql connection, error", zap.String("dnName", dnName), zap.Error(err))
			return result, err
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
			return result, err
		}

		result[dnName] = *db

		logger.Info(ctx, "init single mysql connection success", zap.String("dnName", dnName))
	}

	return result, nil
}

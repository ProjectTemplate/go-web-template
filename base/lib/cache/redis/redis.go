package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

var redisClient map[string]*redis.Client

// GetClient 根据别名获取redis连接，别名必须存在，否则会panic
func GetClient(ctx context.Context, name string) *redis.Client {
	client := redisClient[name]

	if client == nil {
		logger.Error(ctx, "Get Redis Client, client is nil", zap.String("name", name))
		panic("Get Redis Client, client is nil, name: " + name)
	}

	return client
}

// Init 初始化 redis client
//
// 配置文件示例如下：
// 前缀 redis 标识是redis配置
// test 和 test1 是连接的别名，GetClient方法使用别名获取对应的连接
// [redis.test]   #连接一，名字为 test
//
//	addr = "127.0.0.1:6379" #连接地址
//	password = ""           #密码
//	db = 0                  #数据库
//
// [redis.test1] #连接二，名字为 test1
//
//	addr = "127.0.0.1:6379" #连接地址
//	password = ""           #密码
//	db = 1                  #数据库
func Init(ctx context.Context, redisConfigs map[string]config.Redis) {

	if len(redisConfigs) == 0 {
		logger.Error(ctx, "Init Redis, redis config is empty", zap.Any("redisConfigs", redisConfigs))
		panic("Init Redis, redis config is empty")
	}

	redisClient = make(map[string]*redis.Client, len(redisConfigs))
	for name, conf := range redisConfigs {
		logger.Info(ctx, "Init Redis", zap.String("name", name), zap.Any("config", conf))
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		})
		redisClient[name] = client
	}
}

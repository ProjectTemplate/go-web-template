package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
)

var redisClient map[string]*redis.Client

// GetClient 获取 redis client
func GetClient(ctx context.Context, name string) *redis.Client {
	client := redisClient[name]

	if client == nil {
		logger.Error(ctx, "Get Redis Client, client is nil", zap.String("name", name))
		panic("Get Redis Client, client is nil, name: " + name)
	}

	return client
}

// Init 初始化 redis client
func Init(ctx context.Context, redisConfigs map[string]config.Redis) {

	if len(redisConfigs) == 0 {
		logger.SErrorF(ctx, "Init Redis, redis config is empty, configs:+%v", redisConfigs)
		panic("redis config is empty")
	}

	redisClient = make(map[string]*redis.Client, len(redisConfigs))
	for name, conf := range redisConfigs {
		logger.SInfoF(ctx, "Init Redis, name:%s, config:%+v", name, conf)
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		})
		redisClient[name] = client
	}
}

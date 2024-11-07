package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

func Init(ctx context.Context, redisConfigs map[string]config.Redis) map[string]*redis.Client {
	if len(redisConfigs) == 0 {
		logger.SErrorF(ctx, "Init Redis, redis config is empty, configs:+%v", redisConfigs)
		panic("redis config is empty")
	}

	result := make(map[string]*redis.Client, len(redisConfigs))
	for name, conf := range redisConfigs {
		logger.SInfoF(ctx, "Init Redis, name:%s, config:%+v", name, conf)
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		})
		result[name] = client
	}

	return result
}

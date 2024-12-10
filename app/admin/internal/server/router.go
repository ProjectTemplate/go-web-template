package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-web-template/app/admin/internal/api"
	innerRedis "go-web-template/base/lib/cache/redis"
	"go-web-template/base/lib/config"
)

var redisTest *redis.Client
var pingPong *api.PingPong
var postApi *api.PostApi

func InitDependence(ctx context.Context, config *config.Configs) {
	//初始化
	innerRedis.Init(ctx, config.Redis)

	//基础组件
	redisTest = innerRedis.GetClient(ctx, "test")

	//biz

	//api
	pingPong = api.NewPingPong()
	postApi = api.NewPostApi()
}

func RegisterRouter(e *gin.Engine) {
	e.GET("ping", pingPong.Ping)
	e.POST("form", postApi.Form)
	e.POST("json", postApi.Json)
}

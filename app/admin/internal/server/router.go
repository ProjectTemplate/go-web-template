package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/api"
	"go-web-template/base/lib/config"
)

var pingPong *api.PingPongApi
var postApi *api.PostApi
var invokeApi *api.InvokeApi
var getApi *api.GetApi

func InitDependence(ctx context.Context, config *config.Configs) {

	//biz

	//api
	pingPong = api.NewPingPongApi()
	postApi = api.NewPostApi()
	invokeApi = api.NewInvokeApi()
	getApi = api.NewGetApi()
}

func RegisterRouter(e *gin.Engine) {
	e.GET("ping", pingPong.Ping)

	e.GET("get", getApi.Get)

	e.GET("invoke", invokeApi.Invoke)

	e.POST("form", postApi.Form)
	e.POST("json", postApi.Json)
}

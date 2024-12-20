package main

import (
	"context"
	"flag"
	"fmt"
	middleware2 "go-web-template/base/lib/gin/middleware"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-web-template/app/admin/internal/global"
	"go-web-template/app/admin/internal/server"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
	"go-web-template/base/lib/signal"
)

// confFile 配置文件路径
var confFile string

// initCommandLineFlag 初始化命令行参数
func initCommandLineFlag() {
	flag.StringVar(&confFile, "conf", "../../configs/config_dev.toml", "config path, eg: -conf config.yaml")
	flag.Parse()

	absConfPath, err := filepath.Abs(confFile)
	utils.PanicAndPrintIfNotNil(err)

	fmt.Println("confPath: ", confFile)
	fmt.Println("confPath abs: ", absConfPath)
}

func main() {

	ctx := context.Background()
	ctx = utils.WithStartTime(ctx, time.Now())

	initCommandLineFlag()

	config.Init(confFile, global.Configs)

	logger.Init(global.Configs.App.Name, global.Configs.LoggerConfig)
	loggerFlushError := logger.Flush()

	logger.Info(ctx, "start server.", zap.String("confFile", confFile), zap.Any("configs", global.Configs))

	ginMode := gin.ReleaseMode
	if global.Configs.Server.Debug {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	r := gin.New()
	// 中间件处理
	panicRecover := middleware2.PanicRecover(response.AdminInternalErrorReason)
	initContext := middleware2.InitContext(global.Configs.App.Name, response.AdminInternalErrorReason)
	r.Use(panicRecover, initContext)

	//初始化依赖
	server.InitDependence(ctx, global.Configs)
	//注册路由
	server.RegisterRouter(r)

	//启动服务
	go func() {
		address := ":8080"
		if global.Configs != nil && global.Configs.Server.Address != "" {
			address = global.Configs.Server.Address
		}
		err := r.Run(address)
		if err != nil {
			logger.Error(ctx, "server run error", zap.Error(err))
		}
		utils.PanicAndPrintIfNotNil(err)
	}()

	//监听停止信号
	signal.HandleSignal(ctx, func() {
		err := logger.Flush()
		if err != nil && err.Error() != loggerFlushError.Error() {
			logger.Error(ctx, "logger flush error", zap.Error(err))
		}
	})
}

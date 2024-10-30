package main

import (
	"context"
	"flag"
	"fmt"
	"go-web-template/base/lib/middleware"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-web-template/app/admin/internal/global"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

// confFile 配置文件路径
var confFile string

func main() {

	background := context.Background()

	initCommandLineFlag()

	config.Init(confFile, global.Configs)

	logger.Init("go-web", global.Configs.LoggerConfig)

	logger.Info(background, "start server.", zap.String("confFile", confFile), zap.Any("configs", global.Configs))

	ginMode := gin.ReleaseMode
	if global.Configs.Server.Debug {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	r := gin.New()
	// 中间件处理

	r.Use(middleware.PanicRecover())

	r.GET("/ping", func(c *gin.Context) {
		logger.Info(c.Request.Context(), "ping")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("server panic")
	})

	address := ":8080"
	if global.Configs != nil && global.Configs.Server.Address != "" {
		address = global.Configs.Server.Address
	}

	err := r.Run(address)
	if err != nil {
		logger.Error(background, "server run error", zap.Error(err))
	}
	utils.PanicAndPrintIfNotNil(err)

	logger.Info(background, "server run success", zap.String("address", address))
}

// initCommandLineFlag 初始化命令行参数
func initCommandLineFlag() {
	flag.StringVar(&confFile, "conf", "../../configs/config_dev.toml", "config path, eg: -conf config.yaml")
	flag.Parse()

	absConfPath, err := filepath.Abs(confFile)
	utils.PanicAndPrintIfNotNil(err)

	fmt.Println("confPath: ", confFile)
	fmt.Println("confPath abs: ", absConfPath)
}

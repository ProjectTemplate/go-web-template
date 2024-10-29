package main

import (
	"flag"
	"fmt"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/global"
)

// confFile 配置文件路径
var confFile string

func main() {

	initCommandLineFlag()

	config.Init(confFile, global.Configs)

	logger.Init(global.Configs.LoggerConfig)

	logger.New().Warnf("Warn")
	logger.New().Infof("Info")
	logger.New().Errorf("Error")

	//todo user gin.ReleaseMode
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	address := ":8080"
	if global.Configs != nil && global.Configs.Server.Address != "" {
		address = global.Configs.Server.Address
	}
	err := r.Run(address)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("server address", address)
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

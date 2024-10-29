package main

import (
	"flag"
	"fmt"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/global"
)

func main() {

	initCommandLineFlag()

	config.Init(global.ConfFile, global.Configs)

	logger.Init()

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
	flag.StringVar(&global.ConfFile, "conf", "../../configs/config.toml", "config path, eg: -conf config.yaml")
	flag.Parse()

	absConfPath, _ := filepath.Abs(global.ConfFile)
	fmt.Println("confPath: ", global.ConfFile)
	fmt.Println("confPath abs: ", absConfPath)
}

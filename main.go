package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go-web-template/global"
)

func main() {

	initCommandLineFlag()

	initConfig(global.ConfFile)

	//初始化日志

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
	flag.StringVar(&global.ConfFile, "conf", "./configs/config_dev.toml", "config path, eg: -conf config.yaml")
	flag.Parse()

	absConfPath, _ := filepath.Abs(global.ConfFile)
	fmt.Println("confPath: ", global.ConfFile)
	fmt.Println("confPath abs: ", absConfPath)
}

// initConfig 初始化配置
func initConfig(configFile string) {
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	viperHookFunc := mapstructure.ComposeDecodeHookFunc(
		// 字符串转时间 1s 1m 1h 1d
		mapstructure.StringToTimeDurationHookFunc(),
		// 字符串转字符串数组 1,2,3 => [1,2,3]
		mapstructure.StringToSliceHookFunc(","),
	)

	err = viper.Unmarshal(global.Configs, viper.DecodeHook(viperHookFunc))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("config info, config path:%#v, config:%#v\n", configFile, global.Configs)
}

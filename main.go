package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"go-web-template/global"
)

// confPath 配置文件路径
var confPath string

func main() {
	initCommandLineFlag()

	// 加载配置文件
	viper.SetConfigFile(confPath)

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
	fmt.Printf("config info, config path:%#v, config:%#v\n", confPath, global.Configs)

	//todo user gin.ReleaseMode
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	defaultAddress := "0.0.0.0"
	defaultPort := 8888
	if global.Configs != nil && global.Configs.Server.Address != "" {
		defaultAddress = global.Configs.Server.Address
	}
	if global.Configs != nil && global.Configs.Server.Port != 0 {
		defaultPort = global.Configs.Server.Port
	}
	addr := defaultAddress + ":" + fmt.Sprintf("%d", defaultPort)
	fmt.Println(addr)

	err = r.Run(addr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// initCommandLineFlag 初始化命令行参数
func initCommandLineFlag() {
	flag.StringVar(&confPath, "conf", "./configs/config_dev.toml", "config path, eg: -conf config.yaml")
	flag.Parse()

	abs, _ := filepath.Abs(confPath)
	fmt.Println("confPath: ", confPath)
	fmt.Println("confPath abs: ", abs)
}

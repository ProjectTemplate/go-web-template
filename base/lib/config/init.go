package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Init 初始化配置
func Init(configFile string, configStruct any) {
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

	err = viper.Unmarshal(configStruct, viper.DecodeHook(viperHookFunc))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("config info, config path:%#v, config:%#v\n", configFile, configStruct)
}

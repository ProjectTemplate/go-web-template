package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Configs 配置信息
type Configs struct {
	Server       Server       `mapstructure:"server"`
	LoggerConfig LoggerConfig `mapstructure:"log"`
}

// Server 服务器配置
type Server struct {
	Address string `mapstructure:"address"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	// Level 日志级别 zapcore.Level
	Level string `mapstructure:"level"`
	// Console 是否输出到控制台
	Console bool `mapstructure:"console"`
	// Path 日志文件路径
	Path string `mapstructure:"path"`
	// FileName 日志文件名
	FileName string `mapstructure:"file_name"`

	// MaxSize 日志文件最大大小，单位MB
	MaxSize int `mapstructure:"max_size"`
	// MaxAge 日志文件最大个数
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge 日志文件最大保存天数
	MaxAge int `mapstructure:"max_age"`
}

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

package global

import (
	"go-web-template/base/lib/config"
)

var (
	// ConfFile 配置文件路径
	ConfFile string

	// Configs 全局配置
	Configs = &config.Configs{}
)

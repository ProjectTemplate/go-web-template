package conf

import (
	"go-web-template/base/lib/config"
)

type Configs struct {
	Server config.Server `mapstructure:"server"`
}

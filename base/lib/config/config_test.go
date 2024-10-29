package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConfig(t *testing.T) {
	var conf = &Configs{}
	Init("./configs/config.toml", conf)

	assert.NotNil(t, conf)
	assert.NotNil(t, conf.Server)
	assert.NotNil(t, conf.LoggerConfig)

	assert.Equal(t, "127.0.0.1:8080", conf.Server.Address)
	assert.Equal(t, true, conf.LoggerConfig.Console)
	assert.Equal(t, "DEBUG", conf.LoggerConfig.Level)
}

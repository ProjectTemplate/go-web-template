package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseConfig(t *testing.T) {
	var conf = &Configs{}
	Init("./data/config.toml", conf)

	assert.NotNil(t, conf)
	assert.NotNil(t, conf.Server)
	assert.NotNil(t, conf.LoggerConfig)

	assert.Equal(t, "127.0.0.1:8080", conf.Server.Address)
	assert.Equal(t, true, conf.LoggerConfig.Console)
	assert.Equal(t, "DEBUG", conf.LoggerConfig.Level)
	assert.Equal(t, "/var/log/go-web-template", conf.LoggerConfig.Path)
	assert.Equal(t, "server.log", conf.LoggerConfig.FileName)
	assert.Equal(t, 100, conf.LoggerConfig.MaxSize)
	assert.Equal(t, 30, conf.LoggerConfig.MaxBackups)
	assert.Equal(t, 15, conf.LoggerConfig.MaxAge)
}

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	var conf = &Configs{}
	Init("./data/config.toml", conf)

	assert.NotNil(t, conf)
	assert.NotNil(t, conf.Server)
	assert.NotNil(t, conf.LoggerConfig)

	assert.Equal(t, "127.0.0.1:8080", conf.Server.Address)
	assert.Equal(t, true, conf.Server.Debug)

	assert.Equal(t, true, conf.LoggerConfig.Console)
	assert.Equal(t, "DEBUG", conf.LoggerConfig.Level)
	assert.Equal(t, "./", conf.LoggerConfig.Path)
	assert.Equal(t, "server.log", conf.LoggerConfig.FileName)
	assert.Equal(t, 100, conf.LoggerConfig.MaxSize)
	assert.Equal(t, 30, conf.LoggerConfig.MaxBackups)
	assert.Equal(t, 15, conf.LoggerConfig.MaxAge)

	assert.NotEmpty(t, conf.DB)
	assert.Equal(t, 1, len(conf.DB))
	assert.Equal(t, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local&parseTime=True", conf.DB["test"].DSN[0])
	assert.Equal(t, 50, conf.DB["test"].MaxOpenConnections)
	assert.Equal(t, 25, conf.DB["test"].MaxIdleConnections)
	assert.Equal(t, time.Hour, conf.DB["test"].MaxLifeTime)
	assert.Equal(t, time.Minute*10, conf.DB["test"].MaxIdleTime)
	assert.Equal(t, true, conf.DB["test"].IsLogger)
}

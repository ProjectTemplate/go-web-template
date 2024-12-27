package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

func TestRedis(t *testing.T) {
	configs := &config.Configs{}
	config.Init("./data/config.toml", configs)
	logger.Init("TestRedis", configs.LoggerConfig)

	background := context.Background()
	Init(background, configs.Redis)

	c := GetClient(background, "test")
	ping := c.Ping(background)

	assert.Nil(t, ping.Err())
	assert.Equal(t, "PONG", ping.Val())

	c1 := GetClient(background, "test1")
	c1.Set(background, "testKey", "0000000000000000000000000000000000000000000000000000000000000000", time.Second*10)
	usage := c1.MemoryUsage(background, "testKey")
	assert.Nil(t, usage.Err())
	assert.True(t, usage.Val() > 0)
}

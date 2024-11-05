package nacos

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
	"time"
)

func TestNacosCache(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacosCache", configStruct.LoggerConfig)

	background := context.Background()
	client := Init(background, configStruct.Nacos)
	configClient := client.configClients["test"]

	cache := NewConfigCache(configClient)

	group := "test_cache"
	dataId := "number"
	cache.InitConfig(background, group, dataId, UnmarshalToNumber)

	go func() {
		for {
			numberI := cache.GetConfig(background, group, dataId)
			number, ok := numberI.(int)
			assert.True(t, ok)
			logger.SInfoF(background, "number:%v", number)
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 10)
}

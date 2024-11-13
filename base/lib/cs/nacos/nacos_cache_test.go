package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"strconv"
	"testing"
	"time"
)

func TestNacosCache(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacosCache", configStruct.LoggerConfig)

	background := context.Background()
	Init(background, configStruct.Nacos)
	configClient := GetConfigClient(background, "test")

	cache := NewConfigCache(configClient)

	group := "test_cache"
	dataId := "number"
	cache.InitConfig(background, group, dataId, UnmarshalToNumber)

	go func() {
		for i := 0; i < 10; i++ {
			publishConfig, err := configClient.PublishConfig(vo.ConfigParam{
				DataId:  dataId,
				Group:   group,
				Content: strconv.Itoa(i),
			})
			logger.SInfoF(background, "set number:%d", i)
			assert.Nil(t, err)
			assert.True(t, publishConfig)
			time.Sleep(time.Millisecond * 300)
		}
	}()

	go func() {
		for {
			numberI := cache.GetConfig(background, group, dataId)
			number, ok := numberI.(int)
			assert.True(t, ok)
			logger.SInfoF(background, "number:%d", number)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	time.Sleep(time.Second * 10)
	numberI := cache.GetConfig(background, group, dataId)
	number, ok := numberI.(int)
	assert.True(t, ok)
	assert.Equal(t, 9, number)
}

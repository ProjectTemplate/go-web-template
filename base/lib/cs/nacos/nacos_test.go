package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
)

func TestConnection(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)

	background := context.Background()
	Init(background, configStruct.Nacos)

	configClient := GetConfigClient(background, "test")

	configData, err := configClient.GetConfig(vo.ConfigParam{
		Group:  "test",
		DataId: "test",
	})

	assert.Nil(t, err)
	assert.Equal(t, "test", configData)

	namingClient := GetNamingClient(background, "test1")
	instance, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: "test",
		Ip:          "127.0.0.1",
		Port:        8080,
		GroupName:   "test1",
	})

	assert.Nil(t, err)
	assert.True(t, instance)
}

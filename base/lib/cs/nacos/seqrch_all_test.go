package nacos

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
)

func TestSearchAll(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)

	// Test SearchConfigInGroup
	background := context.Background()
	client := Init(background, configStruct.Nacos)
	configClient := client.configClients["test"]

	configItems, err := SearchConfigInGroup(background, "test", configClient)

	assert.Nil(t, err)

	assert.NotEmpty(t, configItems)
}

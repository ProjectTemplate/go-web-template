package nacos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

func TestSearchAll(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)

	// Test SearchConfigInGroup
	background := context.Background()
	Init(background, configStruct.Nacos)
	configClient := GetConfigClient(background, "test")

	configItems, err := SearchConfigInGroup(background, "test", configClient)

	assert.Nil(t, err)

	assert.NotEmpty(t, configItems)
}

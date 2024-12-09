package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
)

func TestGet(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)
	Init(configStruct.FastHttp)

	background := context.Background()
	result := &response{}
	err := Get(background, "http://127.0.0.1:8080/ping", nil, nil, result)

	assert.Nil(t, err)
}

type response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Message string `json:"message"`
	} `json:"data"`
	TraceId string `json:"traceId"`
}

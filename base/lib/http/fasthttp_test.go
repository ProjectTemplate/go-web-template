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

	type request struct {
		Name    string   `url:"name"`
		Age     int      `url:"age"`
		Friends []string `url:"friends"`
	}

	background := context.Background()
	result := &response{}
	params := request{
		Name:    "test",
		Age:     18,
		Friends: []string{"name1", "name2"},
	}

	err := Get(background, "http://127.0.0.1:8080/ping", params, nil, result)

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

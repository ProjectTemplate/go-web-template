package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
	"time"
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

	type response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Message string `json:"message"`
		} `json:"data"`
		TraceId string `json:"traceId"`
	}

	background := context.Background()
	result := &response{}
	params := request{
		Name:    "test",
		Age:     18,
		Friends: []string{"name1", "name2"},
	}

	err := GetTimeOut(background, "http://127.0.0.1:8080/ping", params, nil, time.Second, result)

	assert.Nil(t, err)
}

package http

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-web-template/base/common/constant"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
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

func TestPostForm(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)
	Init(configStruct.FastHttp)

	type request struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Friends []string `json:"friends"`
	}

	type response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Success string `json:"success"`
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

	headers := make(map[string]string)
	headers[constant.HeaderKeyContextType] = constant.ContentTypeForm

	err := PostTimeOut(background, "http://127.0.0.1:8080/form", params, headers, time.Second, result)

	assert.Nil(t, err)
}

func TestPostJson(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestNacos", configStruct.LoggerConfig)
	Init(configStruct.FastHttp)

	type request struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Friends []string `json:"friends"`
	}

	type response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Success string `json:"success"`
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

	headers := make(map[string]string)
	headers[constant.HeaderKeyContextType] = constant.ContentTypeJson
	err := PostTimeOut(background, "http://127.0.0.1:8080/json", params, nil, time.Minute, result)

	assert.Nil(t, err)
}

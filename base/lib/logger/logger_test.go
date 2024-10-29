package logger

import (
	"go-web-template/base/lib/config"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initLogger() {
	configs := config.Configs{}
	config.Init("./data/config.toml", &configs)

	Init(configs.LoggerConfig)
}

func TestName(t *testing.T) {
	initLogger()

	logger := New()
	logger.Info("Info")
	logger.Error("Error")
	defer logger.Sync()
	assert.NotEmpty(t, logger, "日志信息不能为空")
}

func TestMultiSingle(t *testing.T) {
	initLogger()

	logger := New()
	times := 1024
	for i := 0; i < times; i++ {
		logger.Infow("测试打印日志", "name", "name")
	}
}

func TestMultiOpen(t *testing.T) {
	initLogger()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	go func() {
		logger := New()
		times := 1024
		for i := 0; i < times; i++ {
			logger.Infow("1111111111", "name", "name")
		}
		waitGroup.Done()
	}()

	go func() {
		logger := New()
		times := 1024
		for i := 0; i < times; i++ {
			logger.Infow("2222222222", "name", "name")
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

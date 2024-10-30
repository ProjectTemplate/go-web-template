package logger

import (
	"context"
	"go-web-template/base/lib/config"
	"sync"
	"testing"
)

func initLogger() {
	configs := config.Configs{}
	config.Init("./data/config.toml", &configs)

	Init(configs.LoggerConfig)
}

func TestName(t *testing.T) {
	initLogger()

	SInfoF(context.Background(), "Info")
	SErrorF(context.Background(), "Error")
}

func TestMultiSingle(t *testing.T) {
	initLogger()

	times := 1024
	for i := 0; i < times; i++ {
		SInfoF(context.Background(), "测试打印日志", "name", "name")
	}
}

func TestMultiOpen(t *testing.T) {
	initLogger()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			SInfoF(context.Background(), "1111111111", "name", "name")
		}
		waitGroup.Done()
	}()

	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			SInfoF(context.Background(), "2222222222", "name", "name")
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

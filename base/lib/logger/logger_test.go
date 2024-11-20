package logger

import (
	"context"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go.uber.org/zap"
	"sync"
	"testing"
)

func initLogger(name string) {
	configs := config.Configs{}
	config.Init("./data/config.toml", &configs)

	Init(name, configs.LoggerConfig)
}

func TestFBackground(t *testing.T) {
	initLogger("TestF")
	defer Flush()

	ctx := context.Background()

	SDebugF(ctx, "Debug")
	SDebugF(ctx, "Debug, number:%d", 1)

	SInfoF(ctx, "Info")
	SInfoF(ctx, "Info, number:%d", 1)

	SWarnF(ctx, "Warn")
	SWarnF(ctx, "Warn, number:%d", 1)

	SErrorF(ctx, "Error")
	SErrorF(ctx, "Error, number:%d", 1)
}

func TestF(t *testing.T) {
	initLogger("TestF")
	defer Flush()

	ctx := initContext()

	SDebugF(ctx, "Debug")
	SDebugF(ctx, "Debug, number:%d", 1)

	SInfoF(ctx, "Info")
	SInfoF(ctx, "Info, number:%d", 1)

	SWarnF(ctx, "Warn")
	SWarnF(ctx, "Warn, number:%d", 1)

	SErrorF(ctx, "Error")
	SErrorF(ctx, "Error, number:%d", 1)
}

func TestWBackground(t *testing.T) {
	initLogger("TestF")
	defer Flush()

	ctx := initContext()

	Debug(ctx, "Debug")
	Debug(ctx, "Debug", zap.String("number", "1"))

	Info(ctx, "Info")
	Info(ctx, "Info", zap.String("number", "1"))

	Warn(ctx, "Warn")
	Warn(ctx, "Warn", zap.String("number", "1"))

	Error(ctx, "Error")
	Error(ctx, "Error", zap.String("number", "1"))
}

func TestW(t *testing.T) {
	initLogger("TestF")
	defer Flush()

	ctx := initContext()

	Debug(ctx, "Debug")
	Debug(ctx, "Debug", zap.String("number", "1"))

	Info(ctx, "Info")
	Info(ctx, "Info", zap.String("number", "1"))

	Warn(ctx, "Warn")
	Warn(ctx, "Warn", zap.String("number", "1"))

	Error(ctx, "Error")
	Error(ctx, "Error", zap.String("number", "1"))
}

func initContext() context.Context {
	ctx := context.Background()
	defer Flush()

	ctx = utils.WithDomain(ctx, "www.baidu.com")
	ctx = utils.WithTraceId(ctx, "trace-1231231232")
	ctx = utils.WithSpan(ctx, "1")
	return ctx
}

func TestMultiSingle(t *testing.T) {
	initLogger("TestMultiSingle")
	defer Flush()

	times := 1024
	for i := 0; i < times; i++ {
		SInfoF(context.Background(), "测试打印日志,%s,%s", "name", "name")
	}
}

func TestMultiOpen(t *testing.T) {
	initLogger("TestMultiOpen")
	defer Flush()

	ctx := initContext()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			SInfoF(ctx, "1111111111,%s,%s", "name", "name")
		}
		waitGroup.Done()
	}()

	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			SInfoF(ctx, "2222222222,%s,%s", "name", "name")
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

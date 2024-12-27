package logger

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
)

func initLogger(name string) {
	configs := config.Configs{}
	config.Init("./data/config.toml", &configs)

	Init(name, configs.LoggerConfig)
}

func TestWBackground(t *testing.T) {
	initLogger("TestF")

	ctx := initContext()

	Debug(ctx, "Debug")
	Debug(ctx, "Debug", WithHttpField(ctx, zap.String("number", "1"))...)

	Info(ctx, "Info")
	Info(ctx, "Info", WithHttpField(ctx, zap.String("number", "1"))...)

	Warn(ctx, "Warn")
	Warn(ctx, "Warn", WithHttpField(ctx, zap.String("number", "1"))...)

	Error(ctx, "Error")
	Error(ctx, "Error", WithHttpField(ctx, zap.String("number", "1"))...)

}

func TestW(t *testing.T) {
	initLogger("TestF")

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

	ctx = utils.WithHost(ctx, "www.baidu.com")
	ctx = utils.WithTraceId(ctx, "trace-1231231232")
	ctx = utils.WithSpan(ctx, utils.NewSpan("", ""))
	ctx = utils.WithURL(ctx, "hello?name=world")
	ctx = utils.WithStartTime(ctx, time.Now())
	return ctx
}

func TestMultiSingle(t *testing.T) {
	initLogger("TestMultiSingle")

	times := 1024
	for i := 0; i < times; i++ {
		Info(context.Background(), "测试打印日志", zap.String("name", "name"))
	}
}

func TestMultiOpen(t *testing.T) {
	initLogger("TestMultiOpen")

	ctx := initContext()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			Info(ctx, "1111111111", zap.String("name", "name"))
		}
		waitGroup.Done()
	}()

	go func() {
		times := 102400
		for i := 0; i < times; i++ {
			Info(ctx, "2222222222", zap.String("name", "name"))
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()

}

// BenchmarkLogger-8   	  374352	      3036 ns/op
func BenchmarkLogger(b *testing.B) {
	initLogger("TestMultiSingle")

	for i := 0; i < b.N; i++ {
		Info(context.Background(), "测试打印日志", zap.String("name", "name"))
	}
}

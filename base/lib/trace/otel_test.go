package trace

import (
	"context"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// jaeger 服务端测试通过
func TestOtel(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestOtelTrace", configStruct.LoggerConfig)

	Init(context.Background(), configStruct.Otel.Trace)

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			go run()
		}
	}()

	time.Sleep(time.Second * 10000)
}

func run() {
	ctx1, span1 := GetTracer().Start(context.Background(), "span1", trace.WithNewRoot())
	defer span1.End()

	//模拟 http 请求添加header
	newRequest, _ := http.NewRequestWithContext(ctx1, "GET", "127.0.0.1：/ping", nil)
	// 会在header里面添加trace,内容如下 {"Traceparent":"00-e3a781899cd60f96d1ed9d077b103fbd-01a39796f552dee6-01"}
	otel.GetTextMapPropagator().Inject(ctx1, propagation.HeaderCarrier(newRequest.Header))
	logger.Info(ctx1, "inject success", zap.Any("header", newRequest.Header))

	// 添加日志
	span1.AddEvent("message to record what happened in the span, information1")
	span1.AddEvent("message to record what happened in the span, information1")
	span1.SetAttributes(attribute.String("http.method", "get"))
	span1.SetAttributes(attribute.String("http.url", "ping"))
	span1.SetAttributes(attribute.String("user", "123456"))

	i := rand.Int()/1000 + 1
	time.Sleep(time.Millisecond * time.Duration(i))

	mockRequestHttpServer(newRequest)
}

// mockRequestHttp 模拟 http 调用，通过Header传递信息
func mockRequestHttpServer(req *http.Request) {
	// 模拟从Http Header中读取数据
	ctx2 := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(req.Header))
	_, span2 := GetTracer().Start(ctx2, "span2")
	defer span2.End()

	span2.SetAttributes(attribute.String("http.method", "get"))
	span2.SetAttributes(attribute.String("http.url", "ping"))
	span2.SetAttributes(attribute.String("user", "123456"))

	i := rand.Int()/1000 + 1
	time.Sleep(time.Millisecond * time.Duration(i))
}

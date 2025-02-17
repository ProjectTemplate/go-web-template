package trace

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"testing"
	"time"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

// jaeger 服务端测试通过
func TestOtel(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestOtel", configStruct.LoggerConfig)

	Init(context.Background(), configStruct.Otel.Trace)
	go run()

	time.Sleep(time.Second * 10)
}

func run() {
	ctx1, span1 := GetTracer().Start(context.Background(), "span1")
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

	time.Sleep(time.Second)

	//模拟 http 调用，通过Header传递信息
	mockRequestHttp(newRequest)
}

func mockRequestHttp(req *http.Request) {
	// 模拟从Http Header中读取数据
	ctx2 := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(req.Header))
	_, span2 := GetTracer().Start(ctx2, "span2")
	span2.SetAttributes(attribute.String("http.method", "get"))
	span2.SetAttributes(attribute.String("http.url", "ping"))
	span2.SetAttributes(attribute.String("user", "123456"))
	defer span2.End()
	time.Sleep(time.Second)
}

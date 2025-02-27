package metric

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"go.opentelemetry.io/otel/attribute"
	"runtime"
	"testing"
	"time"

	"go.opentelemetry.io/otel/metric"
)

var hostName = utils.GetHostName()

func TestOtelMetric(t *testing.T) {

	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestOtelMetric", configStruct.LoggerConfig)

	ctx := context.Background()

	Init(ctx, configStruct.Otel.Metric)

	meter := NewMeter("meter_test")

	// 定义 goroutine 数量指标
	goroutineCount, err := meter.Int64ObservableGauge("goroutine_count",
		metric.WithDescription("Number of active goroutines"),
	)
	assert.Nil(t, err)

	// 定义内存使用指标
	memAlloc, err := meter.Int64ObservableGauge("memory_allocated",
		metric.WithDescription("Bytes of allocated heap memory"),
		//如果指定 uint 则，上传的数据 key 是 memory_allocated_MB
		metric.WithUnit("MB"),
	)
	assert.Nil(t, err)

	// 注册观测器
	var memStats runtime.MemStats
	obs := []metric.Observable{goroutineCount, memAlloc}

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		runtime.ReadMemStats(&memStats)

		observer.ObserveInt64(goroutineCount, int64(runtime.NumGoroutine()),
			metric.WithAttributes(
				attribute.String("pod", "sass"),
				attribute.String("cluster", configStruct.Otel.Metric.Cluster),
				attribute.String("service_name", configStruct.Otel.Metric.ServiceName),
			))

		observer.ObserveInt64(memAlloc, int64(memStats.Alloc/8/1024/1024))
		return nil
	}, obs...)

	assert.Nil(t, err)

	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			go func() {
				time.Sleep(time.Minute)
			}()
		}
	}()

	time.Sleep(time.Second * 1000)
}

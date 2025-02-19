package metric

import (
	"context"
	"log"
	"runtime"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func TestOtelMetric(t *testing.T) {

	Init(context.Background())

	meter := otel.Meter("service_test")

	// 定义 goroutine 数量指标
	goroutineCount, err := meter.Int64ObservableGauge("goroutine_count_1",
		metric.WithDescription("Number of active goroutines"),
		metric.WithUnit("个"),
	)
	if err != nil {
		log.Fatalf("failed to create goroutine_count gauge: %v", err)
	}

	// 定义内存使用指标
	var memStats runtime.MemStats
	memAlloc, err := meter.Int64ObservableGauge("memory_allocated",
		metric.WithDescription("Bytes of allocated heap memory"),
	)
	if err != nil {
		log.Fatalf("failed to create memory_allocated gauge: %v", err)
	}

	// 定义 CPU 使用率指标（这里简化为获取 goroutine 执行时间）
	cpuTime, err := meter.Int64ObservableGauge("cpu_time",
		metric.WithDescription("CPU time used by goroutines"))

	if err != nil {
		log.Fatalf("failed to create cpu_time gauge: %v", err)
	}

	// 注册观测器
	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			observer.ObserveInt64(goroutineCount, int64(runtime.NumGoroutine()))
			runtime.ReadMemStats(&memStats)
			observer.ObserveInt64(memAlloc, int64(memStats.Alloc))
			observer.ObserveInt64(cpuTime, time.Now().UnixNano())
			return nil
		}, goroutineCount, memAlloc, cpuTime)

	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			go func() {
				time.Sleep(time.Minute)
			}()
		}
	}()

	if err != nil {
		log.Fatalf("failed to register callback: %v", err)
	}

	time.Sleep(time.Second * 100)
}

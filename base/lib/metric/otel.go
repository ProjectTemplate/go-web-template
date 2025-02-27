package metric

import (
	"context"
	"time"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	metric2 "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func Init(ctx context.Context, metricConfig config.Metric) {
	meterProvider, err := newMeterProvider(ctx, metricConfig)
	if err != nil {
		utils.PanicAndPrintIfNotNil(err)
	}

	otel.SetMeterProvider(meterProvider)
}

func NewMeter(name string) metric2.Meter {
	return otel.Meter(name)
}

func newMeterProvider(ctx context.Context, metricConfig config.Metric) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(metricConfig.Endpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(3*time.Second))),
		// resource 为空避免在上报数据添加字段
		metric.WithResource(nil),
	)

	return meterProvider, nil
}

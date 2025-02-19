package metric

import (
	"context"
	"go-web-template/base/common/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"time"
)

func Init(ctx context.Context) {
	meterProvider, err := newMeterProvider(ctx)
	if err != nil {
		utils.PanicAndPrintIfNotNil(err)
	}

	otel.SetMeterProvider(meterProvider)
}

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:64814"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter, metric.WithInterval(3*time.Second)),
		),
	)

	return meterProvider, nil
}

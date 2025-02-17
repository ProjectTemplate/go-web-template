package trace

import (
	"context"
	"go.opentelemetry.io/otel/propagation"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var innerTracer otelTrace.Tracer

func GetTracer() otelTrace.Tracer {
	return innerTracer
}

// Init 初始化Trace
func Init(ctx context.Context, traceConfig config.Trace) {
	if traceConfig.Endpoint == "" || traceConfig.ServiceName == "" {
		utils.PanicAndPrintIfNotNil(errors.New("Init trace config is empty"))
	}

	tracerProvider, err := newTraceProvider(ctx, traceConfig)
	if err != nil {
		utils.PanicAndPrintIfNotNil(err)
		return
	}
	otel.SetTracerProvider(tracerProvider)
	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	innerTracer = otel.Tracer(traceConfig.ServiceName)

	go func() {
		defer utils.RecoverWithFmt()

		select {
		//todo 优雅处理shutdown
		case <-ctx.Done():
			err := tracerProvider.Shutdown(ctx)
			if err != nil {
				logger.Error(ctx, "shutdown otel trace provider failed", zap.Error(err))
				return
			}
			logger.Info(ctx, "shutdown otel trace provider success")
			return
		}
	}()
}

func newTraceProvider(ctx context.Context, traceConfig config.Trace) (*trace.TracerProvider, error) {
	options := make([]otlptracehttp.Option, 0)
	options = append(options, otlptracehttp.WithEndpoint(traceConfig.Endpoint)) // 设置OTEL服务地址
	if traceConfig.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	logger.Info(ctx, "newTraceProvider with options", zap.Any("options", options))
	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	// 创建资源，包含服务元数据
	attributes := make([]attribute.KeyValue, 0)
	attributes = append(attributes, semconv.ServiceNameKey.String(traceConfig.ServiceName))
	if traceConfig.ServiceNamespace != "" {
		attributes = append(attributes, semconv.ServiceNamespaceKey.String(traceConfig.ServiceNamespace))
	}
	if traceConfig.ServiceInstanceID != "" {
		attributes = append(attributes, semconv.ServiceInstanceIDKey.String(traceConfig.ServiceInstanceID))
	}
	if traceConfig.ServiceVersion != "" {
		attributes = append(attributes, semconv.ServiceVersionKey.String(traceConfig.ServiceVersion))
	}

	res, err := resource.New(ctx, resource.WithAttributes(attributes...))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(res))
	return tp, nil
}

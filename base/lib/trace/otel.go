package trace

import (
	"context"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var hostName = utils.GetHostName()
var tracer otelTrace.Tracer

func StartServer(ctx context.Context, spanName string, opts ...otelTrace.SpanStartOption) (context.Context, otelTrace.Span) {
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindServer))
	return GetTracer().Start(ctx, spanName, opts...)
}

func StartInternal(ctx context.Context, spanName string, opts ...otelTrace.SpanStartOption) (context.Context, otelTrace.Span) {
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindInternal))
	return GetTracer().Start(ctx, spanName, opts...)
}

func StartClient(ctx context.Context, spanName string, opts ...otelTrace.SpanStartOption) (context.Context, otelTrace.Span) {
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindClient))
	return GetTracer().Start(ctx, spanName, opts...)
}

func StartProducer(ctx context.Context, spanName string, opts ...otelTrace.SpanStartOption) (context.Context, otelTrace.Span) {
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindProducer))
	return GetTracer().Start(ctx, spanName, opts...)
}

func StartConsumer(ctx context.Context, spanName string, opts ...otelTrace.SpanStartOption) (context.Context, otelTrace.Span) {
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindConsumer))
	return GetTracer().Start(ctx, spanName, opts...)
}

func GetTracer() otelTrace.Tracer {
	if tracer == nil {
		utils.PanicAndPrintIfNotNil(errors.New("Tracer is nil, please init it before use"))
	}
	return tracer
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

	// 设置全局的Provider
	otel.SetTracerProvider(tracerProvider)
	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer = otel.Tracer(traceConfig.ScopeName)

	go func() {
		defer utils.RecoverWithFmt()

		<-ctx.Done()
		err := tracerProvider.Shutdown(ctx)
		if err != nil {
			logger.Error(ctx, "shutdown otel trace provider failed", zap.Error(err))
			return
		}

		logger.Info(ctx, "shutdown otel trace provider success")
	}()
}

func newTraceProvider(ctx context.Context, traceConfig config.Trace) (*trace.TracerProvider, error) {
	options := make([]otlptracehttp.Option, 0)
	options = append(options, otlptracehttp.WithEndpoint(traceConfig.Endpoint))
	if traceConfig.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	logger.Info(ctx, "newTraceProvider new exporter with options", zap.Any("options", options))
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

	attributes = append(attributes, semconv.ServiceInstanceIDKey.String(hostName))

	if traceConfig.ServiceVersion != "" {
		attributes = append(attributes, semconv.ServiceVersionKey.String(traceConfig.ServiceVersion))
	}

	logger.Info(ctx, "newTraceProvider new resource with options", zap.Any("options", options))
	res, err := resource.New(ctx, resource.WithAttributes(attributes...))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(res), trace.WithSampler(&CustomSampler{}))
	return tp, nil
}

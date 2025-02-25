package trace

import (
	"go.opentelemetry.io/otel/sdk/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

var _ trace.Sampler = (*CustomSampler)(nil)

// CustomSampler 自定义采样器，采样会作用在每一个Span上，判断每一个Span是否需要采样
type CustomSampler struct{}

func (a *CustomSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if parameters.Kind == otelTrace.SpanKindServer {
		return trace.SamplingResult{
			Decision:   trace.RecordAndSample,
			Tracestate: otelTrace.SpanContextFromContext(parameters.ParentContext).TraceState(),
		}
	}

	// todo 其它策略

	return trace.SamplingResult{
		Decision:   trace.RecordAndSample,
		Tracestate: otelTrace.SpanContextFromContext(parameters.ParentContext).TraceState(),
	}
}

func (a *CustomSampler) Description() string {
	return "CustomSampler"
}

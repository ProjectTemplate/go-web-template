package logger

import (
	"context"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"go.uber.org/zap"
)

// SpanSuccess 当调用执行成功的时候调用
func SpanSuccess(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(msg, WithSpanField(ctx, fields...)...)
}

// SpanFailed 当调用执行失败的时候调用
func SpanFailed(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(msg, WithSpanField(ctx, fields...)...)
}

// WithSpanField 添加公用的Span日志字段
func WithSpanField(ctx context.Context, fields ...zap.Field) []zap.Field {
	span := utils.GetSpan(ctx)
	span.End()
	//在这里统一结束 Span

	commonFieldCount := 7
	result := make([]zap.Field, commonFieldCount, len(fields)+commonFieldCount)
	result[0] = zap.String(constant.LoggerKeyType, constant.LoggerTypeSpan)
	result[1] = zap.String(constant.LoggerKeyParentSpan, span.GetParentSpan())
	result[2] = zap.String(constant.LoggerKeySpan, span.Span())
	result[3] = zap.String(constant.LoggerKeyName, span.GetName())
	result[4] = zap.Int64(constant.LoggerKeySpanStartUs, span.GetStartTime())
	result[5] = zap.Int64(constant.LoggerKeySpanEndUs, span.GetEndTime())
	result[6] = zap.Int64(constant.LoggerKeySpanDurationUs, span.GetDuration())

	return append(result, fields...)
}
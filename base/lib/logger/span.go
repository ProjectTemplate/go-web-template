package logger

import (
	"context"
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

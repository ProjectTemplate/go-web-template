package logger

import (
	"context"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"go.uber.org/zap"
)

// WithHttpField 添加公用的HTTP日志字段
func WithHttpField(ctx context.Context, fields ...zap.Field) []zap.Field {
	result := make([]zap.Field, 3, len(fields)+3)
	result[0] = zap.String(constant.ContextKeyDomain, utils.GetDomain(ctx))
	result[1] = zap.String(constant.ContextKeyURL, utils.GetURL(ctx))
	result[2] = zap.String(constant.ContextKeyRemoteIp, utils.GetRemoteIP(ctx))

	return append(result, fields...)
}

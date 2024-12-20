package logger

import (
	"context"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"go.uber.org/zap"
)

// WithHttpField 添加公用的HTTP日志字段
func WithHttpField(ctx context.Context, fields ...zap.Field) []zap.Field {
	commonFieldCount := 7
	result := make([]zap.Field, commonFieldCount, len(fields)+commonFieldCount)
	result[0] = zap.String(constant.ContextKeyRemoteIp, utils.GetRemoteIP(ctx))
	result[1] = zap.String(constant.ContextKeyURL, utils.GetURL(ctx))
	result[2] = zap.String(constant.ContextKeyHost, utils.GetHost(ctx))
	result[3] = zap.String(constant.ContextKeyPath, utils.GetPath(ctx))
	result[4] = zap.String(constant.ContextKeyQuery, utils.GetQuery(ctx))
	result[5] = zap.String(constant.ContextKeyPostForm, utils.GetPostForm(ctx))
	result[6] = zap.String(constant.ContextKeyRequestBody, utils.GetRequestBody(ctx))
	return append(result, fields...)
}

package logger

import (
	"context"

	"go.uber.org/zap"

	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
)

// WithHttpField 添加公用的HTTP日志字段
func WithHttpField(ctx context.Context, fields ...zap.Field) []zap.Field {
	commonFieldCount := 7
	result := make([]zap.Field, commonFieldCount, len(fields)+commonFieldCount)
	result[0] = zap.String(constant.LoggerKeyRemoteIp, utils.GetRemoteIP(ctx))
	result[1] = zap.String(constant.LoggerKeyURL, utils.GetURL(ctx))
	result[2] = zap.String(constant.LoggerKeyHost, utils.GetHost(ctx))
	result[3] = zap.String(constant.LoggerKeyPath, utils.GetPath(ctx))
	result[4] = zap.String(constant.LoggerKeyQuery, utils.GetQuery(ctx))
	result[5] = zap.String(constant.LoggerKeyPostForm, utils.GetPostForm(ctx))
	result[6] = zap.String(constant.LoggerKeyRequestBody, utils.GetRequestBody(ctx))
	return append(result, fields...)
}

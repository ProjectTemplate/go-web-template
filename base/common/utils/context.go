package utils

import (
	"context"
	"go-web-template/base/common/constant"
)

func WithTraceId(parent context.Context, traceId string) context.Context {
	return context.WithValue(parent, constant.ContextKeyTraceId, traceId)
}

func GetTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value(constant.ContextKeyTraceId).(string); ok {
		return traceId
	}
	return ""
}

func RpcID(parent context.Context, rpcId string) context.Context {
	return context.WithValue(parent, constant.ContextKeyRpcID, rpcId)
}

func GetRpcID(ctx context.Context) string {
	if rpcId, ok := ctx.Value(constant.ContextKeyRpcID).(string); ok {
		return rpcId
	}
	return ""
}

func WithDomain(parent context.Context, domain string) context.Context {
	return context.WithValue(parent, constant.ContextKeyDomain, domain)
}

func GetDomain(ctx context.Context) string {
	if domain, ok := ctx.Value(constant.ContextKeyDomain).(string); ok {
		return domain
	}
	return ""
}

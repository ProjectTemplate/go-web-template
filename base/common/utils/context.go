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

func WithSpan(parent context.Context, parentSpan string) context.Context {
	return context.WithValue(parent, constant.ContextKeySpan, NewSpan(parentSpan))
}

func GetSpan(ctx context.Context) string {
	if span, ok := ctx.Value(constant.ContextKeySpan).(*Span); ok {
		return span.IncreaseAndGet()
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

func WithURL(parent context.Context, url string) context.Context {
	return context.WithValue(parent, constant.ContextKeyURL, url)
}

func GetURL(ctx context.Context) string {
	if url, ok := ctx.Value(constant.ContextKeyURL).(string); ok {
		return url
	}
	return ""
}

func WithRemoteIP(parent context.Context, remoteIP string) context.Context {
	return context.WithValue(parent, constant.ContextKeyRemoteIp, remoteIP)
}

func GetRemoteIP(ctx context.Context) string {
	if remoteIP, ok := ctx.Value(constant.ContextKeyRemoteIp).(string); ok {
		return remoteIP
	}
	return ""
}

//nolint:staticcheck
package utils

import (
	"context"
	"go-web-template/base/common/constant"
	"time"
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

func WithChildSpan(ctx context.Context, childName string) context.Context {
	child := GetSpan(ctx).Child(childName)
	return WithSpan(ctx, child)
}

func WithSpan(parent context.Context, span *Span) context.Context {
	return context.WithValue(parent, constant.ContextKeySpan, span)
}

func GetSpan(ctx context.Context) *Span {
	if span, ok := ctx.Value(constant.ContextKeySpan).(*Span); ok {
		return span
	}
	return EmptySpan
}

func WithHost(parent context.Context, domain string) context.Context {
	return context.WithValue(parent, constant.ContextKeyHost, domain)
}

func WithPath(parent context.Context, path string) context.Context {
	return context.WithValue(parent, constant.ContextKeyPath, path)
}

func WithQuery(parent context.Context, query string) context.Context {
	return context.WithValue(parent, constant.ContextKeyQuery, query)
}

func WithPostForm(parent context.Context, postForm string) context.Context {
	return context.WithValue(parent, constant.ContextKeyPostForm, postForm)
}

func WithRequestBody(parent context.Context, requestBody string) context.Context {
	return context.WithValue(parent, constant.ContextKeyRequestBody, requestBody)
}

func GetHost(ctx context.Context) string {
	if domain, ok := ctx.Value(constant.ContextKeyHost).(string); ok {
		return domain
	}
	return ""
}

func GetPath(ctx context.Context) string {
	if path, ok := ctx.Value(constant.ContextKeyPath).(string); ok {
		return path
	}
	return ""
}

func GetQuery(ctx context.Context) string {
	if query, ok := ctx.Value(constant.ContextKeyQuery).(string); ok {
		return query
	}
	return ""
}

func GetPostForm(ctx context.Context) string {
	if postForm, ok := ctx.Value(constant.ContextKeyPostForm).(string); ok {
		return postForm
	}
	return ""
}

func GetRequestBody(ctx context.Context) string {
	if requestBody, ok := ctx.Value(constant.ContextKeyRequestBody).(string); ok {
		return requestBody
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

func WithStartTime(parent context.Context, time time.Time) context.Context {
	return context.WithValue(parent, constant.ContextKeyStartTime, time)
}

func GetStartTime(ctx context.Context) time.Time {
	if startTime, ok := ctx.Value(constant.ContextKeyStartTime).(time.Time); ok {
		return startTime
	}
	return time.Now()
}

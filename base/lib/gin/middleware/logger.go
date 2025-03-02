package middleware

import (
	"bytes"
	"io"
	"time"

	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"

	"github.com/gin-gonic/gin"
	otelTrace "go.opentelemetry.io/otel/trace"
)

func InitContext(projectName string, reason response.Reason) gin.HandlerFunc {
	return func(c *gin.Context) {
		spanContext := otelTrace.SpanContextFromContext(c.Request.Context())

		var traceId string
		if spanContext.HasTraceID() {
			traceId = spanContext.TraceID().String()
		} else {
			traceId = c.GetHeader(constant.HeaderKeyTraceId)
		}

		// 默认 traceId
		if traceId == "" {
			traceId = utils.UUID() + "." + projectName
		}

		//span
		parentSpan := c.GetHeader(constant.HeaderKeySpan)
		span := utils.NewSpan(parentSpan, c.Request.URL.Path)

		ctx := utils.WithTraceId(c.Request.Context(), traceId)
		ctx = utils.WithHost(ctx, c.Request.Host)
		ctx = utils.WithPath(ctx, c.Request.URL.Path)
		ctx = utils.WithQuery(ctx, c.Request.URL.Query().Encode())
		ctx = utils.WithPostForm(ctx, c.Request.PostForm.Encode())
		ctx = utils.WithURL(ctx, c.Request.URL.String())
		ctx = utils.WithRemoteIP(ctx, c.RemoteIP())
		ctx = utils.WithSpan(ctx, span)
		ctx = utils.WithStartTime(ctx, time.Now())

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response.ErrorWithReason(c, reason)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		ctx = utils.WithRequestBody(ctx, string(body))

		//设置新的context
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		//打印服务耗时
		logger.SpanSuccess(ctx, "endpoint response time", logger.WithHttpField(ctx)...)
	}
}

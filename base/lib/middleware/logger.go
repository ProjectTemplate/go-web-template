package middleware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
	"time"
)

func InitContext(projectName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// traceId
		traceId := c.GetHeader(constant.HeaderKeyTraceId)
		if traceId == "" {
			traceId = utils.UUID() + "." + projectName
		}

		//span
		parentSpan := c.GetHeader(constant.HeaderKeySpan)

		ctx := utils.WithTraceId(c.Request.Context(), traceId)
		ctx = utils.WithDomain(ctx, c.Request.Host)
		ctx = utils.WithURL(ctx, c.Request.URL.String())
		ctx = utils.WithRemoteIP(ctx, c.RemoteIP())
		ctx = utils.WithSpan(ctx, parentSpan)
		ctx = utils.WithStartTime(ctx, time.Now())

		//设置新的context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

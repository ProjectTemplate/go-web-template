package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
)

func InitContext(projectName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// traceId
		traceId := c.GetHeader(constant.HeaderKeyTraceId)
		if traceId == "" {
			traceId = uuid.New().String() + "." + projectName
		}

		//span
		parentSpan := c.GetHeader(constant.HeaderKeySpan)

		ctx := utils.WithTraceId(c.Request.Context(), traceId)
		ctx = utils.WithDomain(ctx, c.Request.Host)
		ctx = utils.WithURL(ctx, c.Request.URL.String())
		ctx = utils.WithRemoteIP(ctx, c.RemoteIP())
		ctx = utils.WithSpan(ctx, parentSpan)

		//设置新的context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

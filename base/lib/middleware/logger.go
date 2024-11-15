package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
)

func InitContext(projectName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := c.Request.Context()

		//domain
		domain := c.GetHeader(constant.HeaderKeyDomain)
		context = utils.WithDomain(context, domain)

		// traceId
		traceId := c.GetHeader(constant.HeaderKeyTraceId)
		if traceId == "" {
			traceId = uuid.New().String() + "." + projectName
		}
		context = utils.WithTraceId(context, traceId)

		//span
		parentSpan := c.GetHeader(constant.HeaderKeySpan)
		context = utils.WithSpan(context, parentSpan)

		//设置新的context
		c.Request.WithContext(context)

		c.Next()
	}
}

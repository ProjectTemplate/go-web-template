package middleware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
)

// PanicRecover panic recover middleware
func PanicRecover(reason response.Reason) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		if err != nil {
			logger.Error(c.Request.Context(), "CustomRecovery", zap.Any("err", err))
			response.ErrorWithReason(c, reason)
			c.Abort()
			return
		}
		c.Next()
	})
}

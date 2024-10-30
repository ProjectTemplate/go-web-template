package middleware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"net/http"
)

// PanicRecover panic recover middleware
func PanicRecover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		if err != nil {
			logger.Logger().Error("CustomRecovery", zap.Any("err", err))
			c.JSON(http.StatusInternalServerError, "服务内部错误")
			c.Abort()
			return
		}
		c.Next()
	})
}

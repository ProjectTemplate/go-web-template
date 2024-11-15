package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	response := newResponse(c, data, ReasonSuccess)
	c.JSON(http.StatusOK, response)
}

// ErrorWithReason 带有原因的错误响应
func ErrorWithReason(c *gin.Context, reason Reason) {
	response := newResponse(c, nil, reason)
	c.JSON(http.StatusOK, response)
}

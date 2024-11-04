package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	response := newResponse(data, ReasonSuccess)
	c.JSON(http.StatusOK, response)
}

func ErrorWithReason(c *gin.Context, reason Reason) {
	response := newResponse(nil, reason)
	c.JSON(http.StatusOK, response)
}

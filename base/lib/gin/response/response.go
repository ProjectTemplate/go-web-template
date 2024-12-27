package response

import (
	"github.com/gin-gonic/gin"

	"go-web-template/base/common/utils"
)

// Response is the response struct
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
}

func newResponse(c *gin.Context, data interface{}, reason Reason) Response {
	return Response{
		Code:    reason.Code,
		Message: reason.Message,
		Data:    data,
		TraceId: utils.GetTraceId(c.Request.Context()),
	}
}

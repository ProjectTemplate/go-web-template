package api

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/model"
	"go-web-template/base/lib/gin/response"
)

type PingPongApi struct {
}

func NewPingPongApi() *PingPongApi {
	return &PingPongApi{}
}

func (p *PingPongApi) Ping(ginCtx *gin.Context) {
	request := model.PingPongRequest{}
	if err := ginCtx.ShouldBind(&request); err != nil {
		response.ErrorWithReason(ginCtx, response.NewReason(response.AdminInternalErrorCode))
		return
	}

	pongResponse := model.PingPongResponse{
		Message: "pong",
	}
	response.Success(ginCtx, pongResponse)
}

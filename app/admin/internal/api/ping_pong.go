package api

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/model"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
)

type PingPong struct {
}

func NewPingPong() *PingPong {
	return &PingPong{}
}

func (p *PingPong) Ping(ctx *gin.Context) {
	pongResponse := model.PingPongResponse{
		Message: "pong",
	}
	logger.Info(ctx.Request.Context(), "ping pong")
	response.Success(ctx, pongResponse)
}
package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/model"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
	"time"
)

type PingPong struct {
}

func NewPingPong() *PingPong {
	return &PingPong{}
}

func (p *PingPong) Ping(ginCtx *gin.Context) {
	ctx := ginCtx.Request.Context()
	
	logger.Info(ctx, "ping pong start")
	invokeServiceA(ctx)
	invokeServiceB(ctx)
	logger.Info(ctx, "ping pong end")

	pongResponse := model.PingPongResponse{
		Message: "pong",
	}
	response.Success(ginCtx, pongResponse)
}

func invokeServiceA(ctx context.Context) {
	ctx = utils.WithChildSpan(ctx, "serviceA")
	time.Sleep(time.Millisecond * 10)
	logger.Info(ctx, "serviceA success")
	logger.SpanSuccess(ctx, "success")
}

func invokeServiceB(ctx context.Context) {
	ctx = utils.WithChildSpan(ctx, "serviceB")
	time.Sleep(time.Millisecond * 20)
	logger.SpanFailed(ctx, "failed")
}

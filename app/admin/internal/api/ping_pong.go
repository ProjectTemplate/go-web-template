package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/model"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
	"math/rand"
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
	childCtx := utils.WithChildSpan(ctx, "serviceA")
	time.Sleep(time.Millisecond * time.Duration(rand.Int()%1000))
	logger.Info(childCtx, "serviceA time cost", logger.WithSpanField(childCtx)...)
}

func invokeServiceB(ctx context.Context) {
	childCtx := utils.WithChildSpan(ctx, "serviceB")
	time.Sleep(time.Millisecond * time.Duration(rand.Int()%1000))
	logger.Info(childCtx, "serviceB time cost", logger.WithSpanField(childCtx)...)
}

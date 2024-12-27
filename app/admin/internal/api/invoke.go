package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"go-web-template/app/admin/internal/model"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
)

type InvokeApi struct {
}

func NewInvokeApi() *InvokeApi {
	return &InvokeApi{}
}

func (i *InvokeApi) Invoke(ginCtx *gin.Context) {
	ctx := ginCtx.Request.Context()
	logger.Info(ctx, "invoke start")
	invokeServiceA(ctx)
	invokeServiceB(ctx)
	logger.Info(ctx, "invoke end")

	response.Success(ginCtx, model.InvokeResponse{})
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

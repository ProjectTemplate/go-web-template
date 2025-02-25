package api

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	otelTrace "go.opentelemetry.io/otel/trace"
	"time"

	"go-web-template/app/admin/internal/model"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/gin/response"
	"go-web-template/base/lib/logger"
	"go-web-template/base/lib/trace"

	"github.com/gin-gonic/gin"
)

type TraceApi struct {
}

func NewTraceApi() *TraceApi {
	return &TraceApi{}
}

func (t *TraceApi) Invoke(ginCtx *gin.Context) {
	ctx := ginCtx.Request.Context()

	logger.Info(ctx, "trace invoke start")
	t.invokeServiceA(ctx)
	t.invokeServiceB(ctx)
	logger.Info(ctx, "trace invoke end")

	response.Success(ginCtx, model.InvokeResponse{})
}

func (t *TraceApi) invokeServiceA(ctx context.Context) {
	ctx, span := trace.StartInternal(ctx, "TraceApi_invokeServiceA")
	defer span.End()

	ctx = utils.WithChildSpan(ctx, "serviceA")
	time.Sleep(time.Millisecond * 10)
	logger.Info(ctx, "trace serviceA success")
	logger.SpanSuccess(ctx, "success")
}

func (t *TraceApi) invokeServiceB(ctx context.Context) {
	spanContext := otelTrace.SpanContextFromContext(ctx)
	ctx, span := trace.StartInternal(ctx, "TraceApi_invokeServiceB",
		otelTrace.WithAttributes(attribute.String("level", "important")),
	)
	span.AddLink(otelTrace.Link{SpanContext: spanContext})

	defer span.End()

	ctx = utils.WithChildSpan(ctx, "serviceB")
	time.Sleep(time.Millisecond * 20)
	logger.SpanFailed(ctx, "failed")
}

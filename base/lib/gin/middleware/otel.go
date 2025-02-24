package middleware

import (
	"github.com/gin-gonic/gin"
	"go-web-template/base/lib/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strings"
)

func InjectOtelTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		ctx, span := trace.StartServer(ctx, getServerName(c))
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func getServerName(c *gin.Context) string {
	builder := strings.Builder{}

	builder.WriteString(c.Request.Method)
	builder.WriteString(" ")
	builder.WriteString(c.Request.Host)
	builder.WriteString(c.Request.URL.Path)

	return builder.String()
}

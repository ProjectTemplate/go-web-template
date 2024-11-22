package server

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/api"
)

func RegisterRouter(e *gin.Engine) {
	pingPong := api.NewPingPong()
	e.GET("ping", pingPong.Ping)
}

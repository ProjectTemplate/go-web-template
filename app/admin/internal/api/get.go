package api

import (
	"github.com/gin-gonic/gin"

	"go-web-template/app/admin/internal/model"
	"go-web-template/base/lib/gin/response"
)

type GetApi struct {
}

func NewGetApi() *GetApi {
	return &GetApi{}
}

func (p *GetApi) Get(ginCtx *gin.Context) {
	getRequest := model.GetRequest{}

	err := ginCtx.ShouldBind(&getRequest)
	if err != nil {
		response.ErrorWithReason(ginCtx, response.AdminParamErrorReason)
		return
	}

	response.Success(ginCtx, struct{}{})
}

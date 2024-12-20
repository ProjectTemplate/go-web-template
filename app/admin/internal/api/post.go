package api

import (
	"github.com/gin-gonic/gin"
	"go-web-template/app/admin/internal/model"
	"go-web-template/base/lib/gin/response"
)

type PostApi struct {
}

func NewPostApi() *PostApi {
	return &PostApi{}
}

func (p *PostApi) Form(ginCtx *gin.Context) {
	req := model.PostFormReq{}
	err := ginCtx.ShouldBind(&req)
	if err != nil {
		response.ErrorWithReason(ginCtx, response.AdminParamErrorReason)
		return
	}

	response.Success(ginCtx, model.PostResponse{
		Success: "success",
	})
}

func (p *PostApi) Json(ginCtx *gin.Context) {
	req := model.PostJsonReq{}
	err := ginCtx.ShouldBind(&req)
	if err != nil {
		response.ErrorWithReason(ginCtx, response.AdminParamErrorReason)
		return
	}

	response.Success(ginCtx, model.PostResponse{
		Success: "success",
	})
}

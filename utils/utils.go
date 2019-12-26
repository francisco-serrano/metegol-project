package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/metegol-project/viewmodels"
)

func SetResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, viewmodels.BaseResponse{
		StatusCode: statusCode,
		Data:       data,
	})
}

func SetErrorResponse(ctx *gin.Context, statusCode int, err error) {
	SetResponse(ctx, statusCode, err.Error())
}

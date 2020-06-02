package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
)

const ResponseValuesKey = "duc_response_values_key"

type ResponseValue struct {
	StatusCode int
	Response   gin.H
}

func Fail(ctx *gin.Context, baseError dao.BaseError) {
	code := http.StatusInternalServerError
	if baseError.StatusCode() != 0 {
		code = baseError.StatusCode()
	}

	ctx.JSON(code, gin.H{
		"code":    code,
		"message": baseError.Error(),
	})
	ctx.Abort()
}

func Success(ctx *gin.Context, code int, h gin.H) {
	ctx.Set(ResponseValuesKey, &ResponseValue{
		StatusCode: code,
		Response:   h,
	})

	ctx.JSON(code, h)
}

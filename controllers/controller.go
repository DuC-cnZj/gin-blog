package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
)

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
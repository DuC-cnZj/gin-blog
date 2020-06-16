package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/utils/errors"
	"net/http"
	"time"
)

var (
	redisClient = config.Conn.RedisClient
	dbClient    = config.Conn.DB
	esClient    = config.Conn.EsClient
)

type ResponseValue struct {
	StatusCode int
	Response   interface{}
}

func Fail(ctx *gin.Context, baseError errors.BaseError) {
	code := http.StatusInternalServerError
	if baseError.StatusCode() != 0 {
		code = baseError.StatusCode()
	}

	ctx.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": baseError.Error(),
	})
}

func Success(ctx *gin.Context, code int, h gin.H) {
	ctx.Set(config.ResponseValuesKey, &ResponseValue{
		StatusCode: code,
		Response:   h,
	})

	value, exists := ctx.Get(config.AppStartKey)
	if exists {
		t := value.(time.Time)
		ctx.Writer.Header().Set("X-Request-Timing", time.Since(t).String())
	}

	ctx.JSON(code, h)
}

func SuccessString(ctx *gin.Context, code int, str string) {
	ctx.Set(config.ResponseValuesKey, &ResponseValue{
		StatusCode: code,
		Response:   str,
	})

	value, exists := ctx.Get("app_start_key")
	if exists {
		t := value.(time.Time)
		ctx.Writer.Header().Set("X-Request-Timing", time.Since(t).String())
	}

	ctx.String(code, str)
}

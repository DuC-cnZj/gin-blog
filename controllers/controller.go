package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/utils/errors"
	"log"
	"net/http"
	"time"
)

var (
	dbClient    *gorm.DB
	esClient    *elastic.Client
	redisClient *redis.Client
)

func Init() {
	if config.Conn.DB == nil || config.Conn.EsClient == nil || config.Conn.RedisClient == nil {
		log.Fatal("error init controller conn")
	}
	dbClient = config.Conn.DB
	esClient = config.Conn.EsClient
	redisClient = config.Conn.RedisClient
}

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

package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetQueryIntValueWithDefault(ctx *gin.Context, key string, defaultValue int) int {
	query, b := ctx.GetQuery(key)
	if b {
		value, err := strconv.Atoi(query)
		if err == nil {
			return value
		}
	}

	return defaultValue
}

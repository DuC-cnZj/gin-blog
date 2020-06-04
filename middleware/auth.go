package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/services"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, b := services.GetClaimFromCtx(c);b {
			c.Set("userId", claims.ID)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "认证失败",
		})
	}
}

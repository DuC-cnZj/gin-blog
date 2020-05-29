package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/utils"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		var start int
		if len(h) >= 6 {
			start = 6
		}
		header := strings.TrimSpace(h[start:])
		token, err := jwt.ParseWithClaims(header, &utils.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.App.JwtSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*utils.MyCustomClaims); ok && token.Valid {
				c.Set("userId", claims.ID)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg":  "认证失败",
		})
	}
}

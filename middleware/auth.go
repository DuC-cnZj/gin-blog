package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/services"
	"log"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization")[6:])
		log.Println(header)
		token, err := jwt.ParseWithClaims(header, &services.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.App.JwtSecret), nil
		})

		if  err == nil && token.Valid {
			if claims, ok := token.Claims.(*services.MyCustomClaims); ok && token.Valid {
				c.Set("userId", claims.ID)
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(401, gin.H{
			"code": 401,
			"msg": "认证失败",
		})
	}
}

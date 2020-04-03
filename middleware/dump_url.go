package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

func DumpUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("url is=", c.Request.URL.Path)
		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func DumpUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("url is=", c.Request.URL.Path)
		t := time.Now()
		c.Next()
		since := time.Since(t)
		c.Writer.Header().Add("X-Timing", since.String())
		log.Println(since.String())
	}
}

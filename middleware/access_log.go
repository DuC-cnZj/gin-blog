package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/models/dao"
	"io/ioutil"
	"log"
	"time"
)

func HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var  res []byte
		var code int
		var content string
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		content = string(bodyBytes)

		c.Next()

		closer := c.Request.Response
		log.Println("###########")
		log.Println(closer)
		if closer != nil {
			code = closer.StatusCode
			res, _ = ioutil.ReadAll(closer.Body)
		}
		value := c.Value("userId")
		var utype = "App\\SocialiteUser"
		if value == nil {
			value = 0
			utype = ""
		}
		history := models.History{
			Ip:         c.ClientIP(),
			Url:        c.FullPath(),
			Method:     c.Request.Method,
			StatusCode: code,
			UserAgent:  c.Request.UserAgent(),
			Content:    content,
			Response:   string(res),
			VisitedAt: &models.JSONTime{
				Time: time.Now(),
			},
			UserableId:   value.(int),
			UserableType: utype,
		}

		dao.Dao.CreateHistory(&history)
	}
}

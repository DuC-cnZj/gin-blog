package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/models/dao"
	"io/ioutil"
	"time"
)

// todo 记录响应
func HandleLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			res       []byte
			bodyBytes []byte
			code      int
		)

		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

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
			Content:    string(bodyBytes),
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

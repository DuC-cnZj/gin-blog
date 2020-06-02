package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
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

		response := c.Value(controllers.ResponseValuesKey)
		if response != nil {
			response := response.(*controllers.ResponseValue)
			code = response.StatusCode
			res, _ = json.Marshal(response.Response)
		}
		value := c.Value("userId")
		var utype = "App\\SocialiteUser"
		if value == nil {
			value = 0
			utype = ""
		}
		content := string(bodyBytes)
		if content == "" {
			content = "[]"
		}
		history := models.History{
			Ip:         c.ClientIP(),
			Url:        c.Request.RequestURI,
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

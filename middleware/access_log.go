package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/services"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const QueueSize = 100000
const timeout = 1 * time.Second

var LogQueue = make(chan models.History, QueueSize)
var once = &sync.Once{}

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

		if !ShouldLog(c) {
			return
		}
		response := c.Value(config.ResponseValuesKey)
		if response != nil {
			response := response.(*controllers.ResponseValue)
			code = response.StatusCode
			res, _ = json.Marshal(response.Response)
		}
		value := c.Value("userId")

		var utype = "App\\SocialiteUser"
		if value == nil {
			if claims, b := services.GetClaimFromCtx(c); b {
				value = claims.ID
			} else {
				value = 0
				utype = ""
			}
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

		go PushQueue(history)
	}
}

func ShouldLog(c *gin.Context) bool {
	recordMethods := []string{"POST", "GET", "PUT", "DELETE", "PATCH"}
	for _, v := range recordMethods {
		if c.Request.Method == v {
			return true
		}
	}

	return false
}

func PushQueue(history models.History) {
	select {
	case LogQueue <- history:
		// todo 应该push到mq
	case <-time.After(timeout):
		log.Println("log queue full!!")
	}
}

func HandleQueue(ctx context.Context) {
	log.Println("init HandleQueue")
	for {
		select {
		case history, ok := <-LogQueue:
			if ok {
				history.Create()
			}
		case <-ctx.Done():
			once.Do(func() {
				close(LogQueue)
				log.Println("log queue quit.", len(LogQueue))
			})
			return
		}
	}
}

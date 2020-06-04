package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/models/dao"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const QueueSize = 100000
const timeout = 1 * time.Second

var LogQueue = make(chan models.History, QueueSize)
var EndChan = make(chan struct{})
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

		go PushQueue(history)
	}
}

func PushQueue(history models.History) {
	select {
	case LogQueue <- history:
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
				dao.Dao.CreateHistory(&history)
				//log.Println("handle one")
			}
		case <-ctx.Done():
			if len(LogQueue)>0 {
				log.Println("还不能关闭")
				break
			}
			once.Do(func() {
				close(LogQueue)
				close(EndChan)
				log.Println("log queue quit.", len(LogQueue))
			})
			return
		}
	}
}

package middleware

import (
	"bytes"
	"context"

	"github.com/gin-gonic/gin"
	 jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/services"
	"io/ioutil"
	"log"
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
	marshal, err := json.Marshal(history)
	if err != nil {
		return
	}

	publish(marshal)
}

func HandleQueue(name interface{}, ctx context.Context) {
	log.Println("init HandleQueue name: ", name)
	consume(ctx)
	log.Printf("worker :%v ctx done\n", name)
}

// mq handler
var queueName = "access_log"
var json = jsoniter.ConfigCompatibleWithStandardLibrary
var channel *amqp.Channel
var declare amqp.Queue
func Init()  {
	var err error
	channel, err = config.Conn.MQ.Channel()
	if err != nil {
		panic(err)
	}
	declare, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
}

func publish(data []byte) {
	channel.Publish("", declare.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func consume(ctx context.Context) {
	deliveries, err := channel.Consume(declare.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case delivery := <-deliveries:
			var h models.History
			json.Unmarshal(delivery.Body, &h)
			h.Create()
			delivery.Ack(false)
		}
	}
}

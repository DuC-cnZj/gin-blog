package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/models/dao"
	"io/ioutil"
	"log"
	"time"
)

func HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
	//	$data = [
	//		'ip'            => $request->getClientIp(),
	//		'url'           => $request->getPathInfo(),
	//		'method'        => $request->getMethod(),
	//		'content'       => $request->input(),
	//		'user_agent'    => $request->userAgent(),
	//		'visited_at'    => Carbon::now(),
	//		'userable_id'   => $userableId,
	//		'userable_type' => $userableType,
	//		'response'      => $response->getContent(),
	//		'status_code'   => $response->getStatusCode(),
	//];
		c.Next()
		var bytes,res []byte
		var code int
		log.Println(c.Request.Body, c.Request.Response)
		if c.Request.Method == "POST" {
			body := c.Request.Body
			bytes, _ = json.Marshal(&body)
		} else {
			bytes, _ = json.Marshal([]byte{})
		}

		closer := c.Request.Response
		if closer != nil {
			code = closer.StatusCode
			res,_ = ioutil.ReadAll(closer.Body)
			closer.Body.Close()
		}
		value := c.Value("userId")
		var utype string = "App\\SocialiteUser"
		if value ==nil{
			value = 0
			utype = ""
		}
		history := models.History{
			Ip:           c.ClientIP(),
			Url:          c.FullPath(),
			Method:       c.Request.Method,
			StatusCode:   code,
			UserAgent:    c.Request.UserAgent(),
			Content:      string(bytes),
			Response:     string(res),
			VisitedAt:    &models.JSONTime{
				Time: time.Now(),
			},
			UserableId:   value.(int),
			UserableType: utype,
		}

		dao.Dao.CreateHistory(&history)
	}
}
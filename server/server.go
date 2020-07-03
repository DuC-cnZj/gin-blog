package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/streadway/amqp"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models"
	routers "github.com/youngduc/go-blog/routes"
	"log"
	"net/http"
	"sync"
	"time"
)

const DefaultQueueNum = 10

type EmptyWriter struct {
}

func (*EmptyWriter) Write(p []byte) (n int, err error) {
	return
}

type Server struct {
	Config      *config.Config
	DBConn      *gorm.DB
	RedisConn   *redis.Client
	EsConn      *elastic.Client
	MQConn      *amqp.Connection
	HttpServer  *http.Server
	Middlewares gin.HandlersChain
	wg          sync.WaitGroup

	QueueNum int
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) IsDebug() bool {
	return s.Config.App.Debug
}

func (s *Server) IsProduction() bool {
	return !s.Config.App.Debug
}

func (s *Server) Init() {
	e := gin.Default()
	e.Use(s.Middlewares...)

	routers.Init(e)
	controllers.Init()
	models.Init()

	s.HttpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.GetAppConfig().HttpPort),
		Handler:        e,
		ReadTimeout:    s.GetAppConfig().ReadTimeout * time.Second,
		WriteTimeout:   s.GetAppConfig().WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Run() error {
	return s.HttpServer.ListenAndServe()
}

func (s *Server) EnableFastMode() {
	s.SetEmptyLogger()
}

func (s *Server) DisableFastMode(ctx context.Context) {
	s.Middlewares = append(s.Middlewares, func(c *gin.Context) {
		c.Set(config.AppStartKey, time.Now())
	}, middleware.HandleLog())

	num := s.QueueNum
	if num <= 0 {
		num = DefaultQueueNum
	}

	s.wg.Add(num)

	fn := func(name interface{}) {
		defer s.wg.Done()
		middleware.HandleQueue(name, ctx)
	}

	for i := 0; i < num; i++ {
		go fn(i)
	}
}

func (s *Server) SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

func (s *Server) SetEmptyLogger() {
	gin.DefaultWriter = &EmptyWriter{}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) Close() {
	s.wg.Wait()
	log.Println("s.wg.Wait done!")
	s.DBConn.Close()
	log.Println("db close!")
	s.RedisConn.Close()
	log.Println("redis close!")
	s.MQConn.Close()
	log.Println("mq close!")
	log.Println("server close!")
}

func (s *Server) GetAppConfig() *config.App {
	return s.Config.App
}

func (s *Server) GetDBConfig() *config.DB {
	return s.Config.DB
}

func (s *Server) GetESConfig() *config.ES {
	return s.Config.ES
}

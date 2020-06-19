package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models"
	routers "github.com/youngduc/go-blog/routes"
	"net/http"
	"sync"
	"time"
)

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
	HttpServer  *http.Server
	Middlewares gin.HandlersChain
	wg          sync.WaitGroup
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

	s.wg.Add(1)
	go func(ctx context.Context) {
		defer s.wg.Done()
		middleware.HandleQueue(ctx)
	}(ctx)
}

func (s *Server) SetReleaseMode() {
	gin.SetMode(gin.ReleaseMode)
}

func (s *Server) SetEmptyLogger() {
	gin.DefaultWriter = &EmptyWriter{}
}

func (s *Server) Close() {
	s.DBConn.Close()
	s.RedisConn.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.HttpServer.Shutdown(ctx)
	s.wg.Wait()
	s.Close()

	return err
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
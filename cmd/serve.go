package cmd

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/utils/interrupt"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/config"
	routers "github.com/youngduc/go-blog/routes"
)

var (
	configPath string
	fastMode   bool
	server     = &Server{}
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Println(configPath)
		setUp()
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&configPath, "config", "c", ".env", "--config .env")
	serveCmd.PersistentFlags().BoolVarP(&fastMode, "fast", "f", false, "禁用日志和控制台输出")
}

func setUp() {
	if configPath == "" {
		configPath = ".env"
	}
	if !path.IsAbs(configPath) {
		viper.AddConfigPath(".")
	}
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	server.Config = config.Init()
	server.dbConn = config.GetDB()
	server.redisConn = config.GetRedis()
	server.esConn = config.GetElastic()

	controllers.Init()
	models.Init()
}

func run() {
	ctx, done := interrupt.Context()
	defer done()

	if server.IsProduction() {
		server.SetReleaseMode()
	}

	if IsFastMode() {
		server.EnableFastMode()
	} else {
		server.DisableFastMode(ctx)
	}

	server.Init()

	ch := make(chan error)

	go func() {
		<-ctx.Done()

		c, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		err := server.Shutdown(c)

		select {
		case ch <- err:
		default:
		}
	}()

	log.Printf("running in %d....\n", server.GetAppConfig().HttpPort)
	if err := server.Run(); err != nil {
		log.Println("server.Run, err: ", err)

		return
	}

	select {
	case e := <-ch:
		log.Println("异常退出, err: ", e)
	default:
		log.Println("graceful shutdown...")
	}
}

// 急速模式，禁用日志和控制台输出
func IsFastMode() bool {
	return fastMode
}

type EmptyWriter struct {
}

func (*EmptyWriter) Write(p []byte) (n int, err error) {
	return
}

type Server struct {
	Config      *config.Config
	dbConn      *gorm.DB
	redisConn   *redis.Client
	esConn      *elastic.Client
	httpServer  *http.Server
	middlewares gin.HandlersChain
	wg          sync.WaitGroup
}

func (s *Server) IsDebug() bool {
	return s.Config.App.Debug
}

func (s *Server) IsProduction() bool {
	return !s.Config.App.Debug
}

func (s *Server) Init() {
	e := gin.Default()
	e.Use(s.middlewares...)
	routers.Init(e)

	server.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.GetAppConfig().HttpPort),
		Handler:        e,
		ReadTimeout:    s.GetAppConfig().ReadTimeout * time.Second,
		WriteTimeout:   s.GetAppConfig().WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) EnableFastMode() {
	s.SetEmptyLogger()
}

func (s *Server) DisableFastMode(ctx context.Context) {
	s.middlewares = append(s.middlewares, func(c *gin.Context) {
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
	s.dbConn.Close()
	s.redisConn.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
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

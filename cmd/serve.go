package cmd

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
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
}

func run() {
	baseCtx, cancel := context.WithCancel(context.Background())

	if server.IsProduction() {
		server.SetReleaseMode()
	}

	if IsFastMode() {
		server.EnableFastMode()
	} else {
		server.DisableFastMode(baseCtx)
	}

	server.Init()

	go func() {
		log.Printf("running in %d....\n", server.GetAppConfig().HttpPort)
		log.Fatal(server.Run())
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
	<-c
	ctx, cancelFunc := context.WithTimeout(baseCtx, 5*time.Second)
	cancel()

	defer cancelFunc()
	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	<-middleware.EndChan
	server.Close()
	log.Println("平滑关闭")
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

	go func(ctx context.Context) {
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

func (s *Server) GetAppConfig() *config.App {
	return s.Config.App
}

func (s *Server) GetDBConfig() *config.DB {
	return s.Config.DB
}

func (s *Server) GetESConfig() *config.ES {
	return s.Config.ES
}

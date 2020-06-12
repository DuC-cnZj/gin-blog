package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models/dao"
	"github.com/youngduc/go-blog/routes"
	"github.com/youngduc/go-blog/services/oauth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", ".env", "-config .env")
}

func main() {
	flag.Parse()
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
	// 初始化配置
	config.Init()
	dao.Init()
	oauth.Init()

	app := config.Config.App

	e := gin.Default()
	if !config.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	e.Use(middleware.DumpUrl())
	e.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// 初始化路由
	routers.Init(e)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.HttpPort),
		Handler:        e,
		ReadTimeout:    app.ReadTimeout * time.Second,
		WriteTimeout:   app.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx := context.Background()
	go func() {
		log.Println(s.ListenAndServe())
	}()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
	<-c
	err = s.Shutdown(ctx)
	dao.Dao.CloseDB()
	log.Println("平滑关闭")
}

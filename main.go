package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/middleware"
	"github.com/youngduc/go-blog/hello/models/dao"
	"github.com/youngduc/go-blog/hello/routes"
	"github.com/youngduc/go-blog/hello/services/oauth"
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
	gin.SetMode(config.Config.App.RunMode)
	e.Use(middleware.DumpUrl())

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

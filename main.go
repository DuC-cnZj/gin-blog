package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/middleware"
	"github.com/youngduc/go-blog/hello/models/dao"
	"github.com/youngduc/go-blog/hello/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func main() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	// 初始化配置
	config.Init()
	dao.Init()

	app := config.Config.App
	log.Println(config.Config.ES)

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

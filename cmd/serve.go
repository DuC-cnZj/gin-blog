/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/youngduc/go-blog/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/models/dao"
	routers "github.com/youngduc/go-blog/routes"
	"github.com/youngduc/go-blog/services/oauth"
)

var configPath string
var fastMode bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
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
	// 初始化配置
	config.Init()
	oauth.Init()
	dao.Init()
}

const AppStartKey = "app_start_key"

type EmptyWriter struct {
}

func (*EmptyWriter) Write(p []byte) (n int, err error) {
	return
}

func run() {
	baseCtx, cancel := context.WithCancel(context.Background())
	app := config.Config.App

	log.Println("config.Config.App.Debug:", config.Config.App.Debug)
	if !config.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	if IsFastMode() {
		gin.DefaultWriter = &EmptyWriter{}
		log.Println("fastMode")
	}

	e := gin.Default()

	if !IsFastMode() {
		e.Use(func(c *gin.Context) {
			c.Set(AppStartKey, time.Now())
		}, middleware.HandleLog())

		for i := 0; i < 30; i++ {
			go func(ctx context.Context) {
				middleware.HandleQueue(ctx)
			}(baseCtx)
		}
	}


	e.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "x-socket-id"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"X-Request-Timing"},
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

	go func() {
		if config.Config.App.Domain != "" {
			log.Println("autotls running.... ", config.Config.App.Domain)
			log.Fatal(autotls.Run(e, config.Config.App.Domain))
		} else {
			log.Println("gin running....")
			log.Println(s.ListenAndServe())
		}
	}()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
	<-c
	ctx, cancelFunc := context.WithTimeout(baseCtx, 5*time.Second)
	cancel()

	defer cancelFunc()
	err := s.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	<-middleware.EndChan
	dao.Dao.CloseDB()
	log.Println("平滑关闭")
}

// 急速模式，禁用日志和控制台输出
func IsFastMode() bool {
	return fastMode
}

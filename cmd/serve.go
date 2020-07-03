package cmd

import (
	"context"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/server"
	"github.com/youngduc/go-blog/utils/interrupt"
	"log"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/config"
)

var (
	configPath string
	queueNum   int
	fastMode   bool
	srv        *server.Server
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
	serveCmd.PersistentFlags().IntVarP(&queueNum, "qn", "q", 10, "日志处理队列数")
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

	srv = server.NewServer()
	srv.Config = config.Init()
	srv.DBConn = config.GetDB()
	srv.RedisConn = config.GetRedis()
	srv.EsConn = nil
	//srv.EsConn = config.GetElastic()
	srv.MQConn = config.GetMQ()
	srv.QueueNum = queueNum
}

func run() {
	ctx, done := interrupt.Context()
	defer done()

	if srv.IsProduction() {
		srv.SetReleaseMode()
	}
	middleware.Init()

	if IsFastMode() {
		srv.EnableFastMode()
	} else {
		srv.DisableFastMode(ctx)
	}

	srv.Init()

	go func() {
		log.Println(srv.Run())
	}()

	<-ctx.Done()

	c, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := srv.Shutdown(c); err != nil {
		log.Println("异常退出, err: ", err)
	}

	srv.Close()

	log.Println("graceful shutdown...")
}

// 急速模式，禁用日志和控制台输出
func IsFastMode() bool {
	return fastMode
}

package cmd

import (
	"context"
	"errors"
	"github.com/youngduc/go-blog/server"
	"github.com/youngduc/go-blog/utils/interrupt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngduc/go-blog/config"
)

var (
	configPath string
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
	srv.EsConn = config.GetElastic()
}

func run() {
	ctx, done := interrupt.Context()
	defer done()

	if srv.IsProduction() {
		srv.SetReleaseMode()
	}

	if IsFastMode() {
		srv.EnableFastMode()
	} else {
		srv.DisableFastMode(ctx)
	}

	srv.Init()

	ch := make(chan error)

	go func() {
		<-ctx.Done()

		c, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		err := srv.Shutdown(c)

		select {
		case ch <- err:
		default:
		}
	}()

	log.Printf("running in %d....\n", srv.GetAppConfig().HttpPort)
	if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

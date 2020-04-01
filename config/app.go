package config

import (
	"github.com/spf13/viper"
	"time"
)

var Config *config

type config struct {
	App *app
	DB *db
}

func Init() *config {
	Config = &config{
		DB: InitDB(),
		App: InitApp(),
	}

	return Config
}

func InitApp() *app {
	return &app{
		RunMode:      viper.GetString("RUN_MODE"),
		PageSize:     viper.GetInt("PAGE_SIZE"),
		JwtSecret:    viper.GetString("JWT_SECRET"),
		HttpPort:     viper.GetInt("HTTP_PORT"),
		ReadTimeout:  time.Duration(viper.GetInt64("READ_TIMEOUT")),
		WriteTimeout: time.Duration(viper.GetInt64("WRITE_TIMEOUT")),
	}
}

type app struct {
	RunMode string
	PageSize int
	JwtSecret string

	HttpPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

package config

import "github.com/spf13/viper"

type MQ struct {
	UserName string
	Password string
	Host     string
	Port     int
}

func InitMQ() *MQ {
	return &MQ{
		UserName: viper.GetString("MQ_USERNAME"),
		Password: viper.GetString("MQ_PASSWORD"),
		Host:     viper.GetString("MQ_HOST"),
		Port:     viper.GetInt("MQ_PORT"),
	}
}

package config

import "github.com/spf13/viper"

type DB struct {
	Conn     string
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Prefix   string
}

func InitDB() *DB {
	conn := viper.GetString("DB_CONNECTION")
	username := viper.GetString("DB_USERNAME")
	database := viper.GetString("DB_DATABASE")
	pwd := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetInt("DB_PORT")
	prefix := viper.GetString("DB_PREFIX")

	return NewDB(conn, username, pwd, host, database, prefix, port)
}

func NewDB(conn, username, password, host, database, prefix string, port int) *DB {
	return &DB{
		Conn:     conn,
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
		Prefix:   prefix,
	}
}

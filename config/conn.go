package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

var Conn = struct {
	DB          *gorm.DB
	EsClient    *elastic.Client
	RedisClient *redis.Client
	MQ          *amqp.Connection
}{}

func GetRedis() *redis.Client {
	var err error
	Conn.RedisClient = redis.NewClient(Cfg.Redis)

	_, err = Conn.RedisClient.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}

	return Conn.RedisClient
}

func GetDB() *gorm.DB {
	var err error
	dbConfig := Cfg.DB

	Conn.DB, err = gorm.Open(dbConfig.Conn, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Database))

	if err != nil {
		log.Fatal(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbConfig.Prefix + defaultTableName
	}

	Conn.DB.SingularTable(false)
	var logMode bool

	if Cfg.App.Debug {
		logMode = true
	}
	Conn.DB.LogMode(logMode)
	Conn.DB.DB().SetMaxIdleConns(10)
	Conn.DB.DB().SetMaxOpenConns(50)
	Conn.DB.DB().SetConnMaxLifetime(time.Hour)

	return Conn.DB
}

func GetElastic() *elastic.Client {
	// ES
	errorlog := log.New(os.Stdout, "elastic search:", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(
		elastic.SetURL(Cfg.ES.Host),
		elastic.SetErrorLog(errorlog),
		elastic.SetSniff(false),
	)
	// Trace request and response details like this
	// client, err := elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping(Cfg.ES.Host).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(Cfg.ES.Host)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	Conn.EsClient = client

	return client
}

func GetMQ() (conn *amqp.Connection) {
	var (
		err   error
		mqCfg = Cfg.MQ
		url   = fmt.Sprintf("amqp://%s:%s@%s:%d/", mqCfg.UserName, mqCfg.Password, mqCfg.Host, mqCfg.Port)
	)

	conn, err = amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	Conn.MQ = conn

	return
}
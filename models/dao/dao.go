package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/pkg/errors"
	"github.com/youngduc/go-blog/config"
	"log"
	"os"
)

var Dao *dao

type dao struct {
	DB    *gorm.DB
	ES    *elastic.Client
	Redis *redis.Client
}

func Init() {
	var (
		err error
	)

	Dao = &dao{}

	// ES
	errorlog := log.New(os.Stdout, "elastic search:", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(
		elastic.SetURL(config.Config.ES.Host),
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
	info, code, err := client.Ping(config.Config.ES.Host).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion(config.Config.ES.Host)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	Dao.ES = client

	// DB
	dbConfig := config.Config.DB

	Dao.DB, err = gorm.Open(dbConfig.Conn, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
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

	Dao.DB.SingularTable(false)
	var logMode bool

	if config.Config.App.Debug {
		logMode = true
	}
	Dao.DB.LogMode(logMode)
	Dao.DB.DB().SetMaxIdleConns(10)
	Dao.DB.DB().SetMaxOpenConns(100)

	Dao.Redis = redis.NewClient(config.Config.Redis)

	_, err = Dao.Redis.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}
}

func (dao *dao) CloseDB() {
	log.Println("db close")
	defer func() {
		dao.DB.Close()
	}()
}

func (dao *dao) Ping() error {
	var err error
	running := dao.ES.IsRunning()
	if running != true {
		return errors.New("err ES")
	}
	err = dao.DB.DB().Ping()
	return err
}

package dao

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/pkg/errors"
	"github.com/youngduc/go-blog/hello/config"
	"log"
	"os"
)

var Dao *dao

type dao struct {
	db *gorm.DB
	es *elastic.Client
}

func Init() {
	var (
		err error
	)

	Dao = &dao{}

	// es
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(
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
	info, code, err := client.Ping("http://localhost:9200").Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the es version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://localhost:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	Dao.es = client
	//es, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	//Dao.es = es

	// DB
	dbConfig := config.Config.DB

	Dao.db, err = gorm.Open(dbConfig.Conn, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
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

	Dao.db.SingularTable(false)
	Dao.db.LogMode(true)
	Dao.db.DB().SetMaxIdleConns(10)
	Dao.db.DB().SetMaxOpenConns(100)
}

func (dao *dao) CloseDB() {
	defer func() {
		dao.db.Close()
	}()
}

func (dao *dao) Ping() error {
	var err error
	running := dao.es.IsRunning()
	if running != true {
		return errors.New("err es")
	}
	err = dao.db.DB().Ping()
	return err
}

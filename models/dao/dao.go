package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/youngduc/go-blog/hello/config"
	"log"
)

var Dao *dao

type dao struct {
	db *gorm.DB
}

func Init() {
	var (
		err error
	)

	Dao = &dao{}

	dbConfig := config.Config.DB

	Dao.db, err = gorm.Open(dbConfig.Conn, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Database))

	if err != nil {
		log.Fatal(err)
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return dbConfig.Prefix + defaultTableName
	}

	Dao.db.SingularTable(false)
	Dao.db.LogMode(true)
	Dao.db.DB().SetMaxIdleConns(10)
	Dao.db.DB().SetMaxOpenConns(100)
}

func (dao *dao) CloseDB() {
	defer dao.db.Close()
}

func (dao *dao) Ping() error {
	return dao.db.DB().Ping()
}
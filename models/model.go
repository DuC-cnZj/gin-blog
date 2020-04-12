package models

import (
	"database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type Model struct {
	Id int `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime `json:"created_at"`
	UpdatedAt JSONTime `json:"updated_at"`
}

type Paginator struct {
	Total int `json:"total"`
	CurrentPage int `json:"current_page"`
	PerPage int `json:"per_page"`
}

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(b []byte) error  {
	parse, e := time.Parse("2006-01-02 15:04:05", strings.Trim( string(b), "\""))
	fmt.Println(parse,e)
	t.Time = parse
	return nil
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
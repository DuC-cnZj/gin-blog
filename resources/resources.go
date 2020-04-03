package resources

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Author struct {
	Id     int64
	Name   string
	Avatar string
}
type Category struct {
	Id   int64
	Name string
}

type Tag struct {
	Id   int64
	Name string
}
type Comment struct {
	Id        int64
	CommentId int64
	Body      string
	CreatedAt string

	Author  *Author
	Article *Article
}

type Article struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime `json:"created_at"`
	UpdatedAt JSONTime `json:"updated_at"`

	IsTop     bool   `json:"is_top"`
	HeadImage string `json:"head_image"`
	Content   string `json:"content"`
	ContentMd string `json:"content_md"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Display   bool   `json:"display"`

	Author *Author `json:"author"`

	Category *Category `json:"category"`

	Tags []*Tag `json:"tags"`

	Comments []*Comment `json:"comments"`
	//recommendArticles
	//highlight
}

type JSONTime struct {
	time.Time
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
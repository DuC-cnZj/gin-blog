package models

import (
	"database/sql"
	"github.com/youngduc/go-blog/hello/resources"
)

type Article struct {
	Model

	IsTop      bool           `json:"is_top"`
	AuthorId   int64          `json:"author_id"`
	HeadImage  string         `json:"head_image"`
	Content    sql.NullString `json:"content"`
	ContentMd  string         `json:"content_md"`
	Title      string         `json:"title"`
	Desc       string         `json:"desc"`
	Display    bool           `json:"display" gorm:"default:true"`
	CategoryId int64          `json:"category_id"`
	TopAt      *JSONTime      `json:"top_at"`

	Author User `json:"author" gorm:"foreignkey:AuthorId"`

	Category Category `json:"category"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`

	Comments []*Comment `json:"comments"`
	//recommendArticles
	//highlight
}

//$table->increments('id');
//$table->integer('author_id');
//$table->integer('category_id');
//$table->longtext('content')->nullable();
//$table->string('desc');
//$table->string('title');
//$table->string('head_image');
//$table->boolean('display')->default(true);
//$table->timestamp('top_at')->nullable();
//$table->timestamps();
func Get(id int64) *resources.Article {
	var article resources.Article
	//db.Table("articles").Select("id,author_id,category_id,content,`desc`,title,head_image,display,top_at,created_at").
	//	Where("id = ?", id).
	//	First(&article)

	return &article
}

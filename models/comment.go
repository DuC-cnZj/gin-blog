package models

type Comment struct {
	Model

	Visitor string `json:"visitor"`
	Content string `json:"content"`
	ArticleId int64 `json:"article_id"`
	CommentId int64 `json:"comment_id"`
	UserableId int64 `json:"userable_id" gorm:"default:0"`
	UserableType string `json:"userable_type" gorm:"default:''"`

	Author  *User

	Replies []Comment `json:"replies"`
	//Article *Article
}


package models

type Tag struct {
	Model

	Name   string `json:"name"`
	UserId int64  `json:"user_id"`

	Articles []Article `json:"articles" gorm:"many2many:article_tag;association_foreignkey:Id;foreignkey:Id"`
}

package models

type User struct {
	Model

	Name   string `json:"name"`
	Avatar string `json:"avatar"`

	Comments []Comment `gorm:"polymorphic:Userable;"`
}

package models

type Category struct {
	Model

	Name   string `json:"name"`
	UserId int64  `json:"user_id"`
}

package models

import "github.com/youngduc/go-blog/config"

type History struct {
	Model

	Ip           string    `json:"ip"`
	Url          string    `json:"url"`
	Method       string    `json:"method"`
	StatusCode   int       `json:"status_code"`
	UserAgent    string    `json:"user_agent"`
	Address      string    `json:"address"`
	Content      string    `json:"content"`
	Response     string    `json:"response"`
	VisitedAt    *JSONTime `json:"visited_at"`
	UserableId   int       `json:"userable_id"`
	UserableType string    `json:"userable_type"`
}

func (h *History) Create() {
	config.Conn.DB.Create(h)
}

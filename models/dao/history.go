package dao

import (
	"github.com/youngduc/go-blog/models"
)

func (dao *dao) CreateHistory(h *models.History) {
	dao.DB.Create(h)
}

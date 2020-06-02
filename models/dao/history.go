package dao

import (
	"github.com/youngduc/go-blog/models"
	"log"
)

func (dao *dao) CreateHistory(h *models.History) {
	create := dao.DB.Create(h)
	log.Println(create.Error)
}

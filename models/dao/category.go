package dao

import "github.com/youngduc/go-blog/models"

func (dao *dao) IndexCategories() []models.Category {
	var categories []models.Category
	dao.DB.Find(&categories)

	return categories
}

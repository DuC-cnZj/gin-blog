package dao

import "github.com/youngduc/go-blog/hello/models"

func (dao *dao) IndexCategories() []models.Category  {
	var categories []models.Category
	dao.db.Find(&categories)

	return categories
}

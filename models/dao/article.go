package dao

import (
	"github.com/youngduc/go-blog/hello/models"
)

func (dao *dao) ShowArticle(id int) *models.Article {
	article := &models.Article{}
	dao.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments").
		Where("id = ?", id).
		Find(article)

	var UserComments = map[int]*models.Comment{}
	var UserIds []int64
	var users []models.User
	userMap := map[int]models.User{}

	for _, comment := range article.Comments {
		if comment.UserableId != 0 && comment.UserableType == "App\\User"{
			UserComments[comment.Id] = comment
			UserIds = append(UserIds, comment.UserableId)
		}
	}

	dao.db.Where("id in (?)", UserIds).Find(&users)

	for _, v := range users {
		userMap[v.Id] = v
	}
	for _, comment := range UserComments {
		if u,ok := userMap[int(comment.UserableId)];ok {
			comment.Author = &u
		}
	}

	return article
}

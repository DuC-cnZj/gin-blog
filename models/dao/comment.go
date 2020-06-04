package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/utils"
	"log"
	"strconv"
)

func (dao *dao) IndexComments(articleId int) []*models.Comment {
	var comments []*models.Comment
	dao.DB.Where("article_id = ?", articleId).Order("id DESC").Find(&comments)
	log.Println("comments", comments)
	if len(comments) == 0 {
		return nil
	}

	var UserComments = map[int]*models.Comment{}
	var UserIds []int64
	var users []models.User
	userMap := map[int]models.User{}

	for _, comment := range comments {
		if comment.UserableId != 0 && comment.UserableType == "App\\User" {
			UserComments[comment.Id] = comment
			UserIds = append(UserIds, comment.UserableId)
		} else {
			comment.Author.Avatar = ""
			comment.Author.Name = comment.Visitor
		}

		comment.Body = comment.Content
	}

	dao.DB.Where("id in (?)", UserIds).Find(&users)

	for _, v := range users {
		userMap[v.Id] = v
	}
	for _, comment := range UserComments {
		if u, ok := userMap[int(comment.UserableId)]; ok {
			comment.Author = u
		}
	}

	return comments
}

func (dao *dao) StoreComment(c *gin.Context) *models.Comment {
	var reqInfo = struct {
		Content   string `json:"content"`
		CommentId int    `json:"comment_id"`
	}{}

	_ = c.BindJSON(&reqInfo)
	articleId, _ := strconv.Atoi(c.Param("id"))
	content := reqInfo.Content
	commentId := reqInfo.CommentId

	userId, err := utils.ParseUserId(c)
	userType := "App\\SocialiteUser"

	if err == utils.UserIdNotFound {
		userId = 0
		userType = ""
	}

	comment := models.Comment{
		Visitor:      c.ClientIP(),
		Content:      content,
		ArticleId:    int64(articleId),
		CommentId:    int64(commentId),
		UserableId:   int64(userId),
		UserableType: userType,
	}

	comment.Body = comment.Content

	row := new(models.Comment)
	dao.DB.Create(&comment).Scan(&row)

	return &comment
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/utils"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type CommentController struct {
}

func (comment *CommentController) Index(c *gin.Context) {
	var (
		UserComments = map[int]*models.Comment{}
		userMap      = map[int]models.User{}
		UserIds      []int64
		users        []models.User
		comments     []*models.Comment
	)

	param := c.Param("id")
	articleId, _ := strconv.Atoi(param)

	dbClient.Where("article_id = ?", articleId).Order("id DESC").Find(&comments)

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

	dbClient.Where("id in (?)", UserIds).Find(&users)

	for _, v := range users {
		userMap[v.Id] = v
	}
	for _, comment := range UserComments {
		if u, ok := userMap[int(comment.UserableId)]; ok {
			comment.Author = u
		}
	}

	Success(c, http.StatusOK, gin.H{
		"data": comment.recursiveReplies(comments),
	})
}

func (*CommentController) Store(c *gin.Context) {
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
	create := dbClient.Create(&comment)
	if create.Error != nil {
		log.Println(create.Error)
	}
	log.Println(comment)
	create.Scan(&row)

	Success(c, http.StatusCreated, gin.H{
		"data": row,
	})
}

func (*CommentController) recursiveReplies(comments []*models.Comment) interface{} {
	if comments == nil {
		return []*models.Comment{}
	}

	var (
		res = make([]*models.Comment, 0)
		m   = make(map[int]*models.Comment)
	)
	for _, v := range comments {
		m[v.Id] = v
	}

	for _, comment := range m {
		if comment.CommentId == 0 {
			res = append(res, comment)
		} else {
			i, ok := m[int(comment.CommentId)]
			if ok {
				i.Replies = append(i.Replies, comment)
			}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Id > res[j].Id
	})

	return res
}

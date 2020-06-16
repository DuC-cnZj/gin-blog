package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/utils"
	"net/http"
	"sort"
	"strconv"
)

type CommentController struct {
}

func (comment *CommentController) Index(c *gin.Context) {
	const (
		SocialiteUser = models.UserableTypeSocialiteUser
		User          = models.UserableTypeUser
	)

	var (
		comments       []*models.Comment
		users          []models.User
		SocialiteUsers []models.SocialiteUser
		userMap        = map[string][]int64{
			SocialiteUser: {},
			User:          {},
		}
		keyByUserId          = map[int64]models.User{}
		keyBySocialiteUserId = map[int64]models.SocialiteUser{}
	)

	param := c.Param("id")
	articleId, _ := strconv.Atoi(param)

	dbClient.Where("article_id = ?", articleId).Order("id DESC").Find(&comments)

	for _, comment := range comments {
		if comment.UserableId != 0 {
			userMap[comment.UserableType] = append(userMap[comment.UserableType], comment.UserableId)
		} else {
			comment.Author.Avatar = ""
			comment.Author.Name = comment.Visitor
		}

		comment.Body = comment.Content
	}

	dbClient.Where("id in (?)", userMap[User]).Find(&users)
	for _, user := range users {
		keyByUserId[int64(user.Id)] = user
	}
	dbClient.Where("id in (?)", userMap[SocialiteUser]).Find(&SocialiteUsers)
	for _, socialiteUser := range SocialiteUsers {
		keyBySocialiteUserId[int64(socialiteUser.Id)] = socialiteUser
	}

	for _, comment := range comments {
		if comment.UserableId != 0 {
			var ca = models.CommentAuthor{}
			if comment.UserableType == SocialiteUser {
				if user, ok := keyBySocialiteUserId[comment.UserableId]; ok {
					ca.Id = user.Id
					ca.Avatar = user.Avatar
					ca.Name = user.Name
				}
			}

			if comment.UserableType == User {
				if user, ok := keyByUserId[comment.UserableId]; ok {
					ca.Id = user.Id
					ca.Avatar = user.Avatar
					ca.Name = user.Name
				}
			}

			comment.Author = ca
		}
	}

	Success(c, http.StatusOK, gin.H{
		"data": comment.recursiveReplies(comments),
	})
}

func (*CommentController) Store(c *gin.Context) {
	var (
		reqInfo = struct {
			Content   string `json:"content"`
			CommentId int    `json:"comment_id"`
		}{}
		ca       models.CommentAuthor
		su       models.SocialiteUser
		userType = models.UserableTypeSocialiteUser
	)

	_ = c.BindJSON(&reqInfo)
	articleId, _ := strconv.Atoi(c.Param("id"))
	content := reqInfo.Content
	commentId := reqInfo.CommentId

	userId, err := utils.ParseUserId(c)

	if err == utils.UserIdNotFound {
		userId = 0
		userType = ""
	} else {
		if result := dbClient.Where("id = ?", userId).Find(&su); result.Error == nil {
			ca.Id = su.Id
			ca.Name = su.Name
			ca.Avatar = su.Avatar
		}
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

	dbClient.Create(&comment)

	comment.Author = ca

	Success(c, http.StatusCreated, gin.H{
		"data": comment,
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

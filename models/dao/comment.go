package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models"
	"strconv"
)

func (dao *dao) IndexComments(articleId int) []*models.Comment {
	var comments []*models.Comment
	dao.DB.Where("article_id = ?", articleId).Find(&comments)

	var UserComments = map[int]*models.Comment{}
	var UserIds []int64
	var users []models.User
	userMap := map[int]models.User{}

	for _, comment := range comments {
		if comment.UserableId != 0 && comment.UserableType == "App\\User" {
			UserComments[comment.Id] = comment
			UserIds = append(UserIds, comment.UserableId)
		}
	}

	dao.DB.Where("id in (?)", UserIds).Find(&users)

	for _, v := range users {
		userMap[v.Id] = v
	}
	for _, comment := range UserComments {
		if u, ok := userMap[int(comment.UserableId)]; ok {
			comment.Author = &u
		}
	}

	return comments
}

func (dao *dao) StoreComment(c *gin.Context) *models.Comment {
	//	$comment = $article->comments()->create([
	//		'visitor'          => $request->ip(),
	//		'content'          => $htmlContent,
	//		'comment_id'       => $request->input('comment_id', 0),
	//	'userable_id'      => is_null($user) ? 0 : $user->id,
	//		'userable_type'    => is_null($user) ? '' : get_class($user),
	//]);
	content := c.PostForm("content")
	articleId, _ := strconv.Atoi(c.Query("article_id"))
	commentId, _ := strconv.Atoi(c.PostForm("comment_id"))

	userId, err := middleware.ParseUserId(c)
	userType := "App\\SocialiteUser"

	if err == middleware.UserIdNotFound {
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

	row := new(models.Comment)
	dao.DB.Create(&comment).Scan(&row)

	return row
}

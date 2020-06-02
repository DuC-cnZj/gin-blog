package comment_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	param := c.Param("id")
	articleId, _ := strconv.Atoi(param)
	controllers.Success(c, http.StatusOK, gin.H{
		"data": recursiveReplies(dao.Dao.IndexComments(articleId)),
	})
}

func recursiveReplies(comments []*models.Comment) interface{} {
	if comments == nil {
		return []*models.Comment{}
	}
	var res []*models.Comment
	var m = make(map[int]*models.Comment)
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

	return res
}

func Store(c *gin.Context) {
	controllers.Success(c, http.StatusCreated, gin.H{
		"data": dao.Dao.StoreComment(c),
	})
}

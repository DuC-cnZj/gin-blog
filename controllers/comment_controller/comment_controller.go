package comment_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models/dao"
	"net/http"
	"strconv"
)

func Index(c *gin.Context)  {
	param := c.Param("id")
	articleId, e := strconv.Atoi(param)
	if e != nil {
		panic(e)
	}
	c.JSON(http.StatusOK, dao.Dao.IndexComments(articleId))
}

func Store(c *gin.Context)  {
	c.JSON(http.StatusCreated, dao.Dao.StoreComment(c))
}

package category_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models/dao"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, dao.Dao.IndexCategories())
}

package category_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":dao.Dao.IndexCategories(),
	})
}

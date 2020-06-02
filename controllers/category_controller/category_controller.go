package category_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
)

func Index(c *gin.Context) {
	controllers.Success(c, http.StatusOK, gin.H{
		"data":dao.Dao.IndexCategories(),
	})
}

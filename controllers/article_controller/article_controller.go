package article_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models/dao"
	"net/http"
	"strconv"
)

func Index(ctx *gin.Context) {

}

func Show(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	article := dao.Dao.ShowArticle(id)

	ctx.JSON(http.StatusOK, article)
}

func Search(ctx *gin.Context) {

}
func Home(ctx *gin.Context) {

}
func Newest(ctx *gin.Context) {

}
func Popular(ctx *gin.Context) {

}
func Trending(ctx *gin.Context) {

}
func Top(ctx *gin.Context) {

}

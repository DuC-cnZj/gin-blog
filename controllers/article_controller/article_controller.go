package article_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models/dao"
	"log"
	"net/http"
	"strconv"
)

func Index(ctx *gin.Context) {
	page, perPage := 1, 15

	s, b := ctx.GetQuery("page")
	if b {
		i, e := strconv.Atoi(s)
		log.Println(e,i)
		page = i
	}
	s, b = ctx.GetQuery("page_size")
	if b {
		i, _ := strconv.Atoi(s)
		perPage = i
	}
	ctx.JSON(http.StatusOK, dao.Dao.IndexArticles(page, perPage))
}

func Show(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	article := dao.Dao.ShowArticle(id)

	ctx.JSON(http.StatusOK, article)
}

func Search(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.Search(ctx.Query("q")))
}
func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.HomeArticles())
}

func Newest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.NewestArticles())
}
func Popular(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.PopularArticles())
}
func Trending(ctx *gin.Context) {

}
func Top(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.TopArticles())
}

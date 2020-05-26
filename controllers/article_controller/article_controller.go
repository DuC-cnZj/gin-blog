package article_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models/dao"
	"github.com/youngduc/go-blog/services"
	"net/http"
	"strconv"
)

func Index(ctx *gin.Context) {
	page, perPage := 1, 15

	s, b := ctx.GetQuery("page")
	if b {
		i, _ := strconv.Atoi(s)
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

	article, e := dao.Dao.ShowArticle(id)
	if e != nil {
		controllers.Fail(ctx, e)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": article,
	})
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
	var trending services.Trending
	get := trending.Get()
	ctx.JSON(http.StatusOK, dao.Dao.GetArticleByIds(get))
}

func Top(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dao.Dao.TopArticles())
}

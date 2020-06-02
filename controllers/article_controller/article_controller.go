package article_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models/dao"
	"github.com/youngduc/go-blog/services"
	"log"
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

	controllers.Success(ctx, 200, dao.Dao.IndexArticles(page, perPage))
}

func Show(ctx *gin.Context) {
	var trending services.Trending
	id, _ := strconv.Atoi(ctx.Param("id"))

	article, e := dao.Dao.ShowArticle(id)

	if e != nil {
		controllers.Fail(ctx, e)
		return
	}
	trending.Push(article.Id)

	controllers.Success(ctx, 200, gin.H{
		"data": article,
	})
}

func Search(ctx *gin.Context) {
	controllers.Success(ctx, 200, gin.H{
		"data": dao.Dao.Search(ctx.Query("q")),
	})
}

func Home(ctx *gin.Context) {
	controllers.Success(ctx, http.StatusOK, gin.H{
		"data": dao.Dao.HomeArticles(),
	})
}

func Newest(ctx *gin.Context) {
	controllers.Success(ctx, http.StatusOK, gin.H{
		"data": dao.Dao.NewestArticles(),
	})
}

func Popular(ctx *gin.Context) {
	controllers.Success(ctx, http.StatusOK, gin.H{
		"data": dao.Dao.PopularArticles(),
	})
}

func Trending(ctx *gin.Context) {
	var trending services.Trending
	get := trending.Get()
	log.Println("Trending ids", get)
	controllers.Success(ctx, http.StatusOK, gin.H{
		"data": dao.Dao.GetArticleByIds(get),
	})
}

func Top(ctx *gin.Context) {
	controllers.Success(ctx, http.StatusOK, gin.H{
		"data": dao.Dao.TopArticles(),
	})
}

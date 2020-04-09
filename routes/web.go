package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/controllers/article_controller"
	"github.com/youngduc/go-blog/hello/controllers/auth_controller"
	"github.com/youngduc/go-blog/hello/controllers/category_controller"
	"github.com/youngduc/go-blog/hello/controllers/comment_controller"
	"github.com/youngduc/go-blog/hello/models/dao"
	"net/http"
	"runtime"
	"time"
)

func Init(router *gin.Engine) *gin.Engine {
	router.GET("/", Root())

	router.GET("/ping", Ping)

	router.GET("/nav_links", NavLinks)

	router.GET("/login/github", auth_controller.RedirectToProvider)

	router.GET("/login/github/callback", auth_controller.HandleProviderCallback)

	router.GET("/auth/me", auth_controller.Me)

	router.GET("/articles/:id", article_controller.Show)

	router.GET("/articles", article_controller.Index)

	router.GET("/search_articles", article_controller.Search)

	router.GET("/home_articles", article_controller.Home)

	router.GET("/newest_articles", article_controller.Newest)

	router.GET("/popular_articles", article_controller.Popular)

	router.GET("/trending_articles", article_controller.Trending)

	router.GET("/top_articles", article_controller.Top)

	router.GET("/categories", category_controller.Index)

	router.GET("/articles/:id/comments", comment_controller.Index)

	router.POST("/articles", comment_controller.Store)

	return router
}

func NavLinks(context *gin.Context) {
	type Links map[string]string
	context.JSON(http.StatusOK, gin.H{
		"data": []Links{
			{"title": "首页", "link": "/"},
			{"title": "文章", "link": "/articles"},
		},
	})
}

func Root() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.String(http.StatusOK, `
	# welcome! power by %s 
	
	      $$\                    $$\               $$\       $$\                     
	      $$ |                   $  |              $$ |      $$ |                    
	 $$$$$$$ |$$\   $$\  $$$$$$$\\_/$$$$$$$\       $$$$$$$\  $$ | $$$$$$\   $$$$$$\  
	$$  __$$ |$$ |  $$ |$$  _____| $$  _____|      $$  __$$\ $$ |$$  __$$\ $$  __$$\ 
	$$ /  $$ |$$ |  $$ |$$ /       \$$$$$$\        $$ |  $$ |$$ |$$ /  $$ |$$ /  $$ |
	$$ |  $$ |$$ |  $$ |$$ |        \____$$\       $$ |  $$ |$$ |$$ |  $$ |$$ |  $$ |
	\$$$$$$$ |\$$$$$$  |\$$$$$$$\  $$$$$$$  |      $$$$$$$  |$$ |\$$$$$$  |\$$$$$$$ |
	 \_______| \______/  \_______| \_______/       \_______/ \__| \______/  \____$$ |
	                                                                       $$\   $$ |
	                                                                       \$$$$$$  |
	                                                                        \______/ 
	created by duc@2018-%s.
	`, runtime.Version(), time.Now().Format("2006"))
	}
}

func Ping(c *gin.Context) {
	ping := dao.Dao.Ping()
	if ping == nil {
		c.JSON(200, gin.H{
			"status": "ok",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "bad",
	})
}

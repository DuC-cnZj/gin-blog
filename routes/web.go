package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers/article_controller"
	"github.com/youngduc/go-blog/controllers/auth_controller"
	"github.com/youngduc/go-blog/controllers/category_controller"
	"github.com/youngduc/go-blog/controllers/comment_controller"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models/dao"
	"net/http"
	"runtime"
	"time"
)

func Init(router *gin.Engine) *gin.Engine {
	//router.GET("/debug/pprof/profile", func(context *gin.Context) {
	//	pprof.Profile(context.Writer, context.Request)
	//})
	//router.GET("/debug/pprof/cmdline", pprof.Cmdline)
	//router.GET("/debug/pprof/profile", pprof.Profile)
	//router.GET("/debug/pprof/symbol", pprof.Symbol)
	//router.GET("/debug/pprof/trace", pprof.Trace)
	//done

	router.LoadHTMLGlob("templates/*")

	router.GET("/", Root)

	//done
	router.GET("/ping", Ping)

	//done
	router.GET("/nav_links", NavLinks)

	//done
	router.GET("/categories", category_controller.Index)

	//done
	router.GET("/articles/:id/comments", comment_controller.Index)

	//done
	router.GET("/articles/:id", article_controller.Show)

	//done
	router.GET("/articles", article_controller.Index)

	//done
	router.GET("/home_articles", article_controller.Home)

	//done
	router.GET("/newest_articles", article_controller.Newest)

	//done
	router.GET("/popular_articles", article_controller.Popular)

	//done
	router.GET("/top_articles", article_controller.Top)

	//done
	router.POST("/articles/:id/comments", comment_controller.Store)

	//done
	router.GET("/trending_articles", article_controller.Trending)

	//done
	router.GET("/search_articles", article_controller.Search)

	router.GET("/login/github", auth_controller.RedirectToProvider)

	router.GET("/login/github/callback", auth_controller.HandleProviderCallback)

	router.POST("/me",  middleware.Auth(), auth_controller.Me)

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

func Root(context *gin.Context) {
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

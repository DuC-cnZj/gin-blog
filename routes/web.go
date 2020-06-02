package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/controllers/article_controller"
	"github.com/youngduc/go-blog/controllers/auth_controller"
	"github.com/youngduc/go-blog/controllers/category_controller"
	"github.com/youngduc/go-blog/controllers/comment_controller"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models/dao"
	"html/template"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

func Init(router *gin.Engine) *gin.Engine {
	parse, _ := template.New("oauth.tmpl").Parse(temp)

	router.SetHTMLTemplate(parse)

	router.GET("/debug/pprof/profile", func(context *gin.Context) {
		pprof.Profile(context.Writer, context.Request)
	})

	router.GET("/debug/pprof/cmdline", func(context *gin.Context) {
		pprof.Cmdline(context.Writer, context.Request)
	})

	router.GET("/debug/pprof/symbol", func(context *gin.Context) {
		pprof.Symbol(context.Writer, context.Request)
	})

	router.GET("/debug/pprof/trace", func(context *gin.Context) {
		pprof.Trace(context.Writer, context.Request)
	})

	router.GET("/ping", Ping)

	router.GET("/", Root)

	router.GET("/nav_links", NavLinks)

	router.GET("/login/github", auth_controller.RedirectToProvider)

	router.GET("/login/github/callback", auth_controller.HandleProviderCallback)

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

	router.POST("/articles/:id/comments", comment_controller.Store)

	group := router.Group("/", middleware.Auth())

	routes := group.Use(middleware.Auth())
	{
		routes.POST("/me", auth_controller.Me)
	}

	return router
}

func NavLinks(context *gin.Context) {
	type Links map[string]string
	controllers.Success(context, http.StatusOK, gin.H{
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
		controllers.Success(c, 200, gin.H{
			"status": "ok",
		})
		return
	}

	controllers.Success(c, 200, gin.H{
		"status": "bad",
	})
}

const temp = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>oauth github</title>
</head>
<body>
登陆中...
<script>
    window.onload = function () {
        window.top.opener.postMessage("bearer {{ .token }}", "{{ .domain }}");
        window.top.close();
    }
</script>
</body>
</html>
`

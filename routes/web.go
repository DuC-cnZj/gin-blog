package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers/article_controller"
	"github.com/youngduc/go-blog/controllers/auth_controller"
	"github.com/youngduc/go-blog/controllers/category_controller"
	"github.com/youngduc/go-blog/controllers/comment_controller"
	"github.com/youngduc/go-blog/middleware"
	"github.com/youngduc/go-blog/models/dao"
	"html/template"
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
	parse, _ := template.New("oauth.tmpl").Parse(temp)
	router.SetHTMLTemplate(parse)

	use := router.Use(middleware.HandleLog())
	{
		//done
		use.GET("/ping", Ping)

		use.GET("/", Root)

		//done
		use.GET("/nav_links", NavLinks)

		use.GET("/login/github", auth_controller.RedirectToProvider)

		use.GET("/login/github/callback", auth_controller.HandleProviderCallback)

		//done
		use.GET("/articles/:id", article_controller.Show)

		//done
		use.GET("/articles", article_controller.Index)

		//done
		use.GET("/search_articles", article_controller.Search)

		//done
		use.GET("/home_articles", article_controller.Home)

		//done
		use.GET("/newest_articles", article_controller.Newest)

		//done
		use.GET("/popular_articles", article_controller.Popular)

		//todo
		use.GET("/trending_articles", article_controller.Trending)

		//done
		use.GET("/top_articles", article_controller.Top)

		//done
		use.GET("/categories", category_controller.Index)

		//done
		use.GET("/articles/:id/comments", comment_controller.Index)

		//done
		use.POST("/articles/:id/comments", comment_controller.Store)
	}

	routes := router.Use(middleware.Auth(), middleware.HandleLog())
	{
		routes.POST("/me", auth_controller.Me)
	}

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

const temp  = `
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
        window.opener.postMessage("bearer {{ .token }}", "{{ .domain }}");
        window.close();
    }
</script>
</body>
</html>
`
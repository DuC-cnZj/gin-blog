package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/middleware"
	"html/template"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

func Init(router *gin.Engine) *gin.Engine {
	router.Use(middleware.Cors())

	parse, _ := template.New("oauth.tmpl").Parse(temp)

	router.SetHTMLTemplate(parse)

	router.GET("/debug/pprof/", func(context *gin.Context) {
		pprof.Index(context.Writer, context.Request)
	})

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

	router.GET("/system_info", systemInfo)

	router.GET("/ping", Ping)

	router.GET("/", Root)

	router.GET("/nav_links", NavLinks)

	authController := &controllers.AuthController{}
	{
		router.GET("/login/github", authController.RedirectToProvider)

		router.GET("/login/github/callback", authController.HandleProviderCallback)
	}

	articleController := &controllers.ArticleController{}
	{
		router.GET("/articles/:id", articleController.Show)

		router.GET("/articles", articleController.Index)

		router.GET("/search_articles", articleController.Search)

		router.GET("/home_articles", articleController.Home)

		router.GET("/newest_articles", articleController.Newest)

		router.GET("/popular_articles", articleController.Popular)

		router.GET("/trending_articles", articleController.Trending)

		router.GET("/top_articles", articleController.Top)
	}

	categoryController := &controllers.CategoryController{}
	{
		router.GET("/categories", categoryController.Index)
	}

	commentController := &controllers.CommentController{}
	{
		router.GET("/articles/:id/comments", commentController.Index)

		router.POST("/articles/:id/comments", commentController.Store)
	}

	group := router.Group("/", middleware.Auth())

	routes := group.Use(middleware.Auth())
	{
		routes.POST("/me", authController.Me)
	}

	return router
}

func systemInfo(context *gin.Context) {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)
	osDic := make(map[string]interface{}, 0)
	osDic["goOs"] = runtime.GOOS
	osDic["arch"] = runtime.GOARCH
	osDic["mem"] = runtime.MemProfileRate
	osDic["compiler"] = runtime.Compiler
	osDic["version"] = runtime.Version()
	osDic["numGoroutine"] = runtime.NumGoroutine()

	dis, _ := disk.Usage("/")
	diskTotalGB := int(dis.Total) / GB
	diskFreeGB := int(dis.Free) / GB
	diskDic := make(map[string]interface{}, 0)
	diskDic["total"] = diskTotalGB
	diskDic["free"] = diskFreeGB

	mem, _ := mem.VirtualMemory()
	memUsedMB := int(mem.Used) / GB
	memTotalMB := int(mem.Total) / GB
	memFreeMB := int(mem.Free) / GB
	memUsedPercent := int(mem.UsedPercent)
	memDic := make(map[string]interface{}, 0)
	memDic["total"] = memTotalMB
	memDic["used"] = memUsedMB
	memDic["free"] = memFreeMB
	memDic["usage"] = memUsedPercent

	cpuDic := make(map[string]interface{}, 0)
	cpuDic["cpuNum"], _ = cpu.Counts(false)

	controllers.Success(context, 200, gin.H{
		"data": map[string]interface{}{
			"os":  osDic,
			"dis": diskDic,
			"mem": memDic,
			"cpu": cpuDic,
		},
	})
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
	str := fmt.Sprintf(`
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

	controllers.SuccessString(context, http.StatusOK, str)
}

func Ping(c *gin.Context) {
	var (
		dberr    error
		rediserr error
	)
	running := config.Conn.EsClient.IsRunning()
	dberr = config.Conn.DB.DB().Ping()
	_, rediserr = config.Conn.RedisClient.Ping().Result()

	if running && dberr == nil && rediserr == nil {
		controllers.Success(c, 200, gin.H{
			"status": "ok",
		})
		return
	}

	controllers.Success(c, 200, gin.H{
		"status": "bad",
		"data": map[string]interface{}{
			"redis": rediserr == nil,
			"db":    dberr == nil,
			"es":    running,
		},
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
        window.opener.postMessage("bearer {{ .token }}", "{{ .domain }}");
        window.close();
    }
</script>
</body>
</html>
`

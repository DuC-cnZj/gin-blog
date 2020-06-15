package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
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

	router.GET("/system_info", systemInfo)

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

func systemInfo(context *gin.Context) {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
	)
	osDic := make(map[string]interface{}, 0)
	osDic["goOs"] = runtime.GOOS
	osDic["logQueue"] = len(middleware.LogQueue)
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
        window.opener.postMessage("bearer {{ .token }}", "{{ .domain }}");
        window.close();
    }
</script>
</body>
</html>
`

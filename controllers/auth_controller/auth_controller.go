package auth_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/services/oauth"
	"log"
)

func RedirectToProvider(ctx *gin.Context)  {
  oauth.Redirect(ctx)
}

func HandleProviderCallback(ctx *gin.Context)  {
	var code string
	code = ctx.Query("code")
	log.Println(code)
	oauth.HandleProviderCallback(code)
}
func Me(ctx *gin.Context)  {

}

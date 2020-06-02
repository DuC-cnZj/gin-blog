package auth_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/controllers"
	"github.com/youngduc/go-blog/models/dao"
	"github.com/youngduc/go-blog/services/oauth"
	"log"
)

func RedirectToProvider(ctx *gin.Context) {
	oauth.Redirect(ctx)
}

func HandleProviderCallback(ctx *gin.Context) {
	var code string
	code = ctx.Query("code")
	log.Println(code)
	oauth.HandleProviderCallback(code, ctx)
}

func Me(ctx *gin.Context) {
	id := ctx.Value("userId").(int)
	controllers.Success(ctx, 200, gin.H{
		"data": dao.Dao.FindSocialiteUser(id),
	})
}

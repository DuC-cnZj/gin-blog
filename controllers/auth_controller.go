package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/services/oauth"
)

type AuthController struct {
}

func (*AuthController) RedirectToProvider(ctx *gin.Context) {
	oauth.Redirect(ctx)
}

func (*AuthController) HandleProviderCallback(ctx *gin.Context) {
	var code string
	code = ctx.Query("code")
	oauth.HandleProviderCallback(code, ctx)
}

func (*AuthController) Me(ctx *gin.Context) {
	var su models.SocialiteUser
	id := ctx.Value("userId").(int)
	dbClient.
		Where("id = ?", id).
		Find(&su)

	Success(ctx, 200, gin.H{
		"data": su,
	})
}

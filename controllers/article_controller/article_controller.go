package article_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models"
	"log"
	"net/http"
	"strconv"
)

func Index(ctx *gin.Context)  {

}
func Show(ctx *gin.Context)  {
	i, _ := strconv.Atoi(ctx.Param("id"))

	log.Println(int64(i))
	article := models.Get(int64(i))
	ctx.JSON(http.StatusOK, article)
}
func Search(ctx *gin.Context)  {

}
func Home(ctx *gin.Context)  {

}
func Newest(ctx *gin.Context)  {

}
func Popular(ctx *gin.Context)  {

}
func Trending(ctx *gin.Context)  {

}
func Top(ctx *gin.Context)  {

}

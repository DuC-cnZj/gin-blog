package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/youngduc/go-blog/hello/api/v1"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/models"
)

func Init(r *gin.Engine) *gin.Engine {

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(config.Config.App.RunMode)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": config.Config.App,
			"config": config.Config.DB,
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		ping := models.Ping()
		if ping == nil {
			c.JSON(200, gin.H{
				"status": "ok",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "bad",
		})
	})

	apiv1 := r.Group("/api/v1")
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
	}

	return r
}

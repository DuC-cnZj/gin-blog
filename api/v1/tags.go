package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/models"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	m := make(map[string]interface{})
	m["id"]=1
	tags:=models.GetTags(0, 10, m)
	c.JSON(200, gin.H{
		"tags": tags,
	})
}

//新增文章标签
func AddTag(c *gin.Context) {
}

//修改文章标签
func EditTag(c *gin.Context) {
}

//删除文章标签
func DeleteTag(c *gin.Context) {
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/models"
	"net/http"
)

type CategoryController struct {
}

func (*CategoryController) Index(c *gin.Context) {
	var categories []models.Category
	dbClient.Order("id DESC").Find(&categories)

	Success(c, http.StatusOK, gin.H{
		"data": categories,
	})
}

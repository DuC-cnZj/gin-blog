package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/config"
	"strings"
)

var UserIdNotFound = errors.New("user id not found")

type MyCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func ParseUserId(c *gin.Context) (int, error) {
	value := c.Value("userId")
	if value != nil {
		return value.(int), nil
	}
	getHeader := c.GetHeader("Authorization")
	if getHeader == "" {
		return 0, UserIdNotFound
	}
	header := strings.TrimSpace(getHeader[6:])
	token, err := jwt.ParseWithClaims(header, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.App.JwtSecret), nil
	})
	if err == nil && token.Valid {
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			return claims.ID, nil
		}
	}

	return 0, UserIdNotFound
}

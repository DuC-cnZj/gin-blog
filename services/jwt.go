package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/utils"
	"strings"
	"time"
)

func GenToken(id int) (string, error) {
	c := &utils.MyCustomClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(config.Cfg.App.JwtSecret))

	return ss, err
}

func GetClaimFromCtx(c *gin.Context) (*utils.MyCustomClaims, bool) {
	h := c.GetHeader("Authorization")
	var start int
	if len(h) >= 6 {
		start = 6
	}
	header := strings.TrimSpace(h[start:])
	token, err := jwt.ParseWithClaims(header, &utils.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.App.JwtSecret), nil
	})

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(*utils.MyCustomClaims); ok && token.Valid {
			return claims, true
		}
	}

	return nil, false
}

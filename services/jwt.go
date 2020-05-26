package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/utils"
	"time"
)



func GenToken(id int) (string, error)  {
	c := &utils.MyCustomClaims{
		ID:             id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute*10).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(config.Config.App.JwtSecret))
	fmt.Printf("%v %v", ss, err)
	return ss, err
}

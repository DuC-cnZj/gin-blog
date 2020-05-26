package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/youngduc/go-blog/config"
	"time"
)

type MyCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func GenToken(id int) (string, error)  {
	c := &MyCustomClaims{
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

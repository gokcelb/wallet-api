package auth

import (
	"fmt"
	"time"

	"github.com/gokcelb/wallet-api/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type TokenService struct {
	conf config.JWTConf
}

func NewTokenService(conf config.JWTConf) *TokenService {
	return &TokenService{conf}
}

func (ts *TokenService) Create() (string, error) {
	key := []byte(ts.conf.Secret)

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().
			Add(time.Minute * time.Duration(ts.conf.ValidityDurationInMin)).Unix(),
		IssuedAt: time.Now().Unix(),
		Issuer:   ts.conf.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Error(err)
	}

	return ss, err
}

func (ts *TokenService) Decode(tokenString string, ctx echo.Context) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(ts.conf.Secret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok {
		log.Print("Token claims are good to go")
	} else {
		log.Error("Token claims are just not it", err)
	}

	if !token.Valid {
		log.Error("Token is invalid", err)
	}

	return token, err
}

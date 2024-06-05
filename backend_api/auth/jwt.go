package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret key")
var accessTokenExpiration = (time.Minute * 10)
var refreshTokenExpiration = (time.Hour * 24)

func GenerateLoginTokens(username string) (LoginTokens, error) {
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(accessTokenExpiration).Unix(),
		})
	accesstokenstring, err := accesstoken.SignedString(secretKey)
	if err != nil {
		return LoginTokens{}, err
	}

	refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(refreshTokenExpiration).Unix(),
		})
	refreshtokenstring, err := refreshtoken.SignedString(secretKey)

	if err != nil {
		return LoginTokens{}, err
	}

	return LoginTokens{AccessToken: accesstokenstring, RefreshToken: refreshtokenstring}, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

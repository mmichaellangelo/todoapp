package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte("secret key")
var refreshSecret = []byte("secret key")
var accessTokenExpiration = (time.Second * 10)
var refreshTokenExpiration = (time.Second * 30)

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"userid"`
	jwt.RegisteredClaims
}

func GenerateLoginTokens(username string, userid int64) (LoginTokens, error) {
	access, err := GenerateAccessToken(username, userid)
	if err != nil {
		return LoginTokens{}, err
	}

	refresh, err := GenerateRefreshToken(username, userid)
	if err != nil {
		return LoginTokens{}, err
	}

	return LoginTokens{AccessToken: access, RefreshToken: refresh}, nil
}

func GenerateAccessToken(username string, userid int64) (string, error) {
	accessClaims := &Claims{
		Username: username,
		UserID:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
		},
	}

	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accesstokenstring, err := accesstoken.SignedString(accessSecret)
	if err != nil {
		return "", err
	}

	return accesstokenstring, nil
}

func isTokenValid(token *jwt.Token) error {
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return fmt.Errorf("token is expired")
		}
		return nil
	} else {
		return fmt.Errorf("invalid token")
	}
}

func GenerateRefreshToken(username string, userid int64) (string, error) {
	refreshClaims := &Claims{
		Username: username,
		UserID:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshtokenstring, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", err
	}

	return refreshtokenstring, nil
}

func VerifyAccessToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	if err != nil {
		return err
	}

	err = isTokenValid(token)
	if err != nil {
		return err
	}

	return nil
}

func VerifyRefreshToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil {
		return err
	}

	err = isTokenValid(token)
	if err != nil {
		return err
	}

	return nil
}

func GetClaimsFromToken(tokenString string) (*Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	err = isTokenValid(token)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func RefreshAccess(refresh string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refresh, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	newAccess, err := GenerateAccessToken(claims.Username, claims.UserID)
	if err != nil {
		return "", err
	}

	return newAccess, nil
}

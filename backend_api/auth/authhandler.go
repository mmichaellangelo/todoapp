package auth

import (
	"fmt"
	"mykale/todobackendapi/db"
)

type AuthHandler struct {
	db *db.DBPool
}

func NewAuthHandler(db *db.DBPool) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) AuthenticateUser(userid int64, username string, passwordattempt string, passwordhashed string) (LoginTokens, error) {
	err := CheckPasswordHash(passwordattempt, passwordhashed)
	if err != nil {
		return LoginTokens{}, fmt.Errorf("incorrect password")
	}

	tokens, err := GenerateLoginTokens(username, userid)
	if err != nil {
		return LoginTokens{}, err
	}
	return tokens, nil
}

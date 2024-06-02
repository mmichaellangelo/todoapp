package auth

import (
	"context"
	"fmt"
	"mykale/todobackendapi/db"
)

type AccountAuth struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	db *db.DBPool
}

func NewAuthHandler(db *db.DBPool) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) AuthenticateUser(username string, password string) (token string, err error) {
	auth, err := h.GetAccountAuthByUsername(username)
	if err != nil {
		return "", err
	}

	if auth.Username != username {
		return "", fmt.Errorf("username mismatch")
	}
	if password != auth.Password {
		return "", fmt.Errorf("incorrect password")
	}

	token, err = CreateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *AuthHandler) GetAccountAuthByID(id int64) (AccountAuth, error) {
	var accountauth AccountAuth
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, password FROM account WHERE id=$1", id)

	if err != nil {
		fmt.Println("Error querying account")
		return accountauth, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&accountauth.ID, &accountauth.Username, &accountauth.Password)
		if err != nil {
			fmt.Println("Error scanning rows")
			return accountauth, err
		}
	}

	return accountauth, nil
}

func (h *AuthHandler) GetAccountAuthByUsername(username string) (AccountAuth, error) {
	var accountauth AccountAuth
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, password FROM auth WHERE username=$1", username)

	if err != nil {
		fmt.Println("Error querying account")
		return accountauth, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&accountauth.ID, &accountauth.Username, &accountauth.Password)
		if err != nil {
			fmt.Println("Error scanning rows")
			return accountauth, err
		}
	}

	return accountauth, nil
}

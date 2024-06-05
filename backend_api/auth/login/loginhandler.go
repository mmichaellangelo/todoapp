package login

import (
	"context"
	"encoding/json"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"net/mail"
	"regexp"
	"time"
)

type LoginHandler struct {
	db             *db.DBPool
	authhandler    *auth.AuthHandler
	accounthandler *account.AccountHandler
}

var (
	LoginRE = regexp.MustCompile(`^\/login\/$`)
)

func NewLoginHandler(db *db.DBPool, authhandler *auth.AuthHandler, accounthandler *account.AccountHandler) http.Handler {
	return &LoginHandler{db: db, authhandler: authhandler, accounthandler: accounthandler}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && LoginRE.MatchString(r.URL.Path):
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tokens, err := h.HandleLogin(r, ctx)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		jsonAuthInfo, err := json.Marshal(tokens)

		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonAuthInfo)
	}
}

func (h *LoginHandler) HandleLogin(r *http.Request, ctx context.Context) (auth.LoginTokens, error) {
	r.ParseMultipartForm(0)
	emailorusername := r.FormValue("emailorusername")
	passwordattempt := r.FormValue("password")

	var account account.Account
	email, err := validMailAddress(emailorusername)
	if err != nil {
		account, err = h.accounthandler.GetByUsername(emailorusername)
		if err != nil {
			return auth.LoginTokens{}, err
		}
	} else {
		account, err = h.accounthandler.GetByEmail(email)
		if err != nil {
			return auth.LoginTokens{}, err
		}
	}

	tokens, err := h.authhandler.AuthenticateUser(account.Username, passwordattempt, account.PasswordHashed)
	if err != nil {
		return auth.LoginTokens{}, err
	}
	return tokens, nil
}

func validMailAddress(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}

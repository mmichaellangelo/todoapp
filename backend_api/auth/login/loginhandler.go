package login

import (
	"context"
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
	LoginRE   = regexp.MustCompile(`^\/login\/?$`)
	LogoutRE  = regexp.MustCompile(`^\/logout\/?$`)
	RefreshRE = regexp.MustCompile(`^\/login\/refresh\/?$`)
)

func NewLoginHandler(db *db.DBPool, authhandler *auth.AuthHandler, accounthandler *account.AccountHandler) http.Handler {
	return &LoginHandler{db: db, authhandler: authhandler, accounthandler: accounthandler}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	// Login route
	case r.Method == http.MethodPost && LoginRE.MatchString(r.URL.Path):
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tokens, err := h.HandleLogin(r, ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		err = h.storeRefreshToken(tokens.RefreshToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("error storing refresh token: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("accesstoken", tokens.AccessToken)
		w.Header().Set("refreshtoken", tokens.RefreshToken)

		w.Write([]byte("login success"))
		return
	// logout route
	case r.Method == http.MethodPost && LogoutRE.MatchString(r.URL.Path):
		refresh := r.Header.Get("refreshtoken")
		if refresh == "" {
			http.Error(w, "refresh token not supplied", http.StatusBadRequest)
			return
		}
		err := h.deleteRefreshToken(refresh)
		if err != nil {
			http.Error(w, fmt.Sprintf("error deleting refresh token: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("logout success"))
		return

	// Refresh route
	case r.Method == http.MethodPost && RefreshRE.MatchString(r.URL.Path):
		fmt.Println("REFRESH")
		refresh := r.Header.Get("refreshtoken")
		if refresh == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		access, err := auth.RefreshAccess(refresh)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("accesstoken", access)
		w.WriteHeader(http.StatusAccepted)
		return
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

	tokens, err := h.authhandler.AuthenticateUser(account.ID, account.Username, passwordattempt, account.PasswordHashed)
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

func (h *LoginHandler) storeRefreshToken(token string) error {
	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO refreshtokens (token) 
	VALUES ($1)`, token)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (h *LoginHandler) deleteRefreshToken(token string) error {
	rows, err := h.db.Pool.Query(context.Background(),
		`DELETE FROM refreshtokens 
		WHERE token=$1`, token)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

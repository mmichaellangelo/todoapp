package account

import (
	"context"
	"encoding/json"
	"fmt"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"time"
)

type Account struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHashed string    `json:"password"`
	DateCreated    time.Time `json:"date_created"`
	DateEdited     time.Time `json:"date_edited"`
}

type AccountHandler struct {
	db          *db.DBPool
	authhandler *auth.AuthHandler
}

var (
	AccountRE       = regexp.MustCompile(`^\/accounts\/$`)
	AccountREWithID = regexp.MustCompile(`^\/accounts\/(\d+)\/?$`)
)

// INSTANCE
func NewAccountHandler(db *db.DBPool, authhandler *auth.AuthHandler) *AccountHandler {
	return &AccountHandler{db: db, authhandler: authhandler}
}

// ROUTES
func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	//----------------------------------- GET ALL ACCOUNTS
	case r.Method == http.MethodGet && AccountRE.MatchString(r.URL.Path):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// MUST HAVE API KEY ----------- ** TODO
		accounts, err := h.GetAll()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		resp, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(resp)
		return
		//----------------------------------- GET ACCOUNT BY ID
	case r.Method == http.MethodGet && AccountREWithID.MatchString(r.URL.Path):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		groups := AccountREWithID.FindStringSubmatch(r.URL.Path)
		if len(groups) != 2 {
			w.WriteHeader(400)
			w.Write([]byte("invalid request"))
			return
		}
		id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error parsing integer"))
			return
		}
		account, err := h.GetByID(id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		resp, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(resp)
		return
	// ----------------------------- CREATE ACCOUNT
	case r.Method == http.MethodPost && AccountRE.MatchString(r.URL.Path):
		err := r.ParseMultipartForm(0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		validemail, err := mail.ParseAddress(email)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid email address"))
		}

		account, err := h.Create(username, validemail.Address, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		tokens, err := h.authhandler.AuthenticateUser(account.ID, account.Username, password, account.PasswordHashed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("accesstoken", tokens.AccessToken)
		w.Header().Set("refreshtoken", tokens.RefreshToken)

		w.Write(resp)
		return

	default:
		return
	}
}

// CREATE
func (h *AccountHandler) Create(username string, email string, password_plaintext string) (Account, error) {
	// hash password
	password_hashed, err := auth.HashPassword(password_plaintext)
	if err != nil {
		return Account{}, err
	}
	// query db and insert
	rows, err := h.db.Pool.Query(context.Background(), `INSERT INTO accounts (username, email, password) 
														VALUES ($1, $2, $3) 
														RETURNING id, username, email, password, date_created, date_edited`,
		username, email, password_hashed)
	if err != nil {
		return Account{}, err
	}
	// prepare data to return
	var account Account
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// READ
func (h *AccountHandler) GetAll() ([]Account, error) {
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, password, date_created, date_edited FROM accounts")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (h *AccountHandler) GetByUsername(username string) (Account, error) {
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, password, date_created, date_edited FROM accounts WHERE username=$1", username)

	if err != nil {
		return Account{}, err
	}

	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return account, err
		}
	}

	return account, nil
}

func (h *AccountHandler) GetByID(id int64) (Account, error) {

	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, password, date_created, date_edited FROM accounts WHERE id=$1", id)

	if err != nil {
		return Account{}, err
	}
	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return account, err
		}
	}
	if account.ID == 0 {
		return Account{}, fmt.Errorf("error finding account")
	}

	return account, nil
}

func (h *AccountHandler) GetByEmail(email string) (Account, error) {

	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, password, date_created, date_edited FROM accounts WHERE email=$1", email)

	if err != nil {
		return Account{}, err
	}

	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return account, err
		}
	}

	return account, nil
}

// UPDATE
func (h *AccountHandler) UpdateUsername(id int64, newusername string) (Account, error) {
	rows, err := h.db.Pool.Query(context.Background(), `UPDATE accounts SET username=$1 WHERE id=$2
														RETURNING id, username, email, password, date_created, date_edited`,
		newusername, id)

	if err != nil {
		return Account{}, err
	}

	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return Account{}, err
		}
	}

	return account, nil
}

func (h *AccountHandler) UpdateEmail(id int64, newemail string) (Account, error) {
	rows, err := h.db.Pool.Query(context.Background(), `UPDATE accounts SET email=$1 WHERE id=$2
														RETURNING id, username, email, password, date_created, date_edited`,
		newemail, id)

	if err != nil {
		return Account{}, err
	}
	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return account, err
		}
	}

	return account, nil
}

func (h *AccountHandler) UpdatePassword(id int64, newpassword_plaintext string) (Account, error) {
	// hash password
	newpassword_hashed, err := auth.HashPassword(newpassword_plaintext)
	if err != nil {
		return Account{}, err
	}
	rows, err := h.db.Pool.Query(context.Background(), `UPDATE accounts SET password=$1 WHERE id=$2
														RETURNING id, username, email, password, date_created, date_edited`,
		newpassword_hashed, id)
	if err != nil {
		return Account{}, err
	}

	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return Account{}, err
		}
	}

	return account, nil
}

// DELETE
func (h *AccountHandler) Delete(id int64) (Account, error) {
	rows, err := h.db.Pool.Query(context.Background(), `DELETE FROM accounts 
													 WHERE id=$1, 
													 RETURNING id, username, email, password, date_created, date_edited`,
		id)

	if err != nil {
		return Account{}, err
	}

	defer rows.Close()
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return Account{}, err
		}
	}

	return account, nil
}

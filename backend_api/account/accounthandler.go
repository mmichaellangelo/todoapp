package account

import (
	"context"
	"fmt"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"regexp"
	"time"
)

type AccountHandler struct {
	db *db.DBPool
}

var (
	AccountRE             = regexp.MustCompile(`^\/accounts\/$`)
	AccountREWithID       = regexp.MustCompile(`^\/accounts\/id\/[0-9]+$`)
	AccountREWithUsername = regexp.MustCompile(`^\/accounts\/username\/[A-z0-9-_]+$`)
)

func NewAccountHandler(db *db.DBPool) http.Handler {
	return &AccountHandler{db: db}
}

func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// Get all accounts
	case r.Method == http.MethodGet && AccountRE.MatchString(r.URL.Path):

	}
}

// CREATE
func (h *AccountHandler) CreateAccount(username string, email string, password_plaintext string) (Account, error) {
	// hash password
	password_hashed, err := auth.HashPassword(password_plaintext)
	if err != nil {
		return Account{}, err
	}
	rows, err := h.db.Pool.Query(context.Background(), "INSERT INTO account (username, email, password, date_created) VALUES($1, $2, $3, $4)", username, email, password_hashed, time.Now())
	if err != nil {
		return Account{}, err
	}
	var account Account
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.DateCreated)
		if err != nil {
			return Account{}, err
		}
	}
	return account, nil
}

// READ
func (h *AccountHandler) GetAllAccounts() ([]Account, error) {
	var accounts []Account
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, date_created FROM account")
	if err != nil {
		fmt.Println("Error querying accounts:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.DateCreated)
		if err != nil {
			fmt.Println("Error querying row:", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (h *AccountHandler) GetAccountByUsername(username string) (Account, error) {
	var account Account
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, date_created FROM account WHERE username=$1", username)

	if err != nil {
		fmt.Println("Error querying account")
		return account, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.DateCreated)
		if err != nil {
			fmt.Println("Error scanning rows")
			return account, err
		}
	}

	return account, nil
}

func (h *AccountHandler) GetAccountByID(id int64) (Account, error) {
	var account Account
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, date_created FROM account WHERE id=$1", id)

	if err != nil {
		fmt.Println("Error querying account")
		return account, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.DateCreated)
		if err != nil {
			fmt.Println("Error scanning rows")
			return account, err
		}
	}

	return account, nil
}

func (h *AccountHandler) GetAccountByEmail(email string) (Account, error) {
	var account Account
	rows, err := h.db.Pool.Query(context.Background(), "SELECT id, username, email, date_created FROM account WHERE email=$1", email)

	if err != nil {
		fmt.Println("Error querying account")
		return account, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.DateCreated)
		if err != nil {
			fmt.Println("Error scanning rows")
			return account, err
		}
	}

	return account, nil
}

// UPDATE
// TODO ------------------------------------------------------

// DELETE
func (h *AccountHandler) DeleteAccountByID(id int64) error {
	_, err := h.db.Pool.Query(context.Background(), "DELETE FROM account WHERE id=$1", id)
	return err
}

func (h *AccountHandler) DeleteAccountByUsername(username string) error {
	_, err := h.db.Pool.Query(context.Background(), "DELETE FROM account WHERE username=$1", username)
	return err
}

func (h *AccountHandler) DeleteAccountByEmail(email string) error {
	_, err := h.db.Pool.Query(context.Background(), "DELETE FROM account WHERE email=$1", email)
	return err
}

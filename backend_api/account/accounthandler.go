package account

import (
	"context"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"regexp"
)

type AccountHandler struct {
	db *db.DBPool
}

var (
	AccountRE             = regexp.MustCompile(`^\/accounts\/$`)
	AccountREWithID       = regexp.MustCompile(`^\/accounts\/id\/[0-9]+$`)
	AccountREWithUsername = regexp.MustCompile(`^\/accounts\/username\/[A-z0-9-_]+$`)
)

func NewAccountHandler(db *db.DBPool) *AccountHandler {
	return &AccountHandler{db: db}
}

func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// Get all accounts
	case r.Method == http.MethodGet && AccountRE.MatchString(r.URL.Path):

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
	for rows.Next() {
		err := rows.Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHashed, &account.DateCreated, &account.DateEdited)
		if err != nil {
			return Account{}, err
		}
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

package account

import (
	"mykale/todobackendapi/db"
	"regexp"
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
	db *db.DBPool
}

var (
	AccountRE       = regexp.MustCompile(`^\/accounts\/$`)
	AccountREWithID = regexp.MustCompile(`^\/accounts\/(\d+)\/?$`)

	ListRE       = regexp.MustCompile(`^\/accounts\/lists\/?`)
	ListREWithID = regexp.MustCompile(`^\/accounts\/lists\/(\d+)\/?`)

	TodoRE       = regexp.MustCompile(`^\/accounts\/todos\/?`)
	TodoREWithID = regexp.MustCompile(`^\/accounts\/todos\/(\d+)\/?`)
)

func NewAccountHandler(db *db.DBPool) *AccountHandler {
	return &AccountHandler{db: db}
}

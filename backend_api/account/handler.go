package account

import (
	"context"
	"encoding/json"
	"fmt"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/password"
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
	db *db.DBPool
}

var (
	AccountRE       = regexp.MustCompile(`^\/accounts\/$`)
	AccountREWithID = regexp.MustCompile(`^\/accounts\/(\d+)\/?$`)
)

func NewAccountHandler(db *db.DBPool) *AccountHandler {
	return &AccountHandler{db: db}
}
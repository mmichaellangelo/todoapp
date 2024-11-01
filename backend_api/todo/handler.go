package todo

import (
	"context"
	"database/sql"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type TodoHandler struct {
	db             *db.DBPool
	accounthandler *account.AccountHandler
}

type Todo struct {
	ID             int64         `json:"id"`
	Body           string        `json:"body"`
	List_ID        int64         `json:"list_id"`
	Completed      bool          `json:"completed"`
	Account_ID     int64         `json:"account_id"`
	Date_Created   time.Time     `json:"date_created"`
	Date_Edited    time.Time     `json:"date_edited"`
	Permissions_ID sql.NullInt64 `json:"permissions_id"`
}

var (
	TodoRE       = regexp.MustCompile(`^\/accounts\/(\d+)\/todos\/?$`)
	TodoREWithID = regexp.MustCompile(`^\/todos\/(\d+)\/?$`)
)

func NewTodoHandler(db *db.DBPool, accounthandler *account.AccountHandler) *TodoHandler {
	return &TodoHandler{db: db, accounthandler: accounthandler}
}




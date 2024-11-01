package list

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/permission"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/functions"
	"mykale/todobackendapi/todo"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type ListHandler struct {
	db                *db.DBPool
	accounthandler    *account.AccountHandler
	todohandler       *todo.TodoHandler
	permissionhandler *permission.PermissionHandler
}

type List struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Account_ID     int64     `json:"account_id"`
	Parent_List_ID int64     `json:"parent_list_id"`
	Permissions_ID int64     `json:"permissions_id"`
	Date_Created   time.Time `json:"date_created"`
	Date_Edited    time.Time `json:"date_edited"`
}

type ListWithTodos struct {
	ID             int64       `json:"id"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Account_ID     int64       `json:"account_id"`
	Parent_List_ID int64       `json:"parent_list_id"`
	Permissions_ID int64       `json:"permissions_id"`
	Date_Created   time.Time   `json:"date_created"`
	Date_Edited    time.Time   `json:"date_edited"`
	Todos          []todo.Todo `json:"todos"`
}

var (
	// all lists for user by id
	ListRE = regexp.MustCompile(`^\/accounts\/(\d+)\/lists\/?$`)
	// list by id
	ListREWithID = regexp.MustCompile(`^\/lists\/(\d+)\/?$`)
)

func NewListHandler(db *db.DBPool, accounthandler *account.AccountHandler, todohandler *todo.TodoHandler, permissionhandler *permission.PermissionHandler) *ListHandler {
	return &ListHandler{db: db, accounthandler: accounthandler, todohandler: todohandler, permissionhandler: permissionhandler}
}




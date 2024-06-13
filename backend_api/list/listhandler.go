package list

import (
	"context"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/todo"
	"net/http"
	"time"
)

type ListHandler struct {
	db             *db.DBPool
	accounthandler *account.AccountHandler
	todohandler    *todo.TodoHandler
	authhandler    *auth.AuthHandler
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

func NewListHandler(db *db.DBPool, accounthandler *account.AccountHandler, todohandler *todo.TodoHandler, authhandler *auth.AuthHandler) *ListHandler {
	return &ListHandler{db: db, accounthandler: accounthandler, todohandler: todohandler, authhandler: authhandler}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {

	}
}

func (h *ListHandler) Create(title string, description string, account_id int64, parent_list_id int64, permissions_id int64) (List, error) {
	// if -1 provided for permissions_id, create a new permission and use that
	if permissions_id == -1 {
		permission, err := h.authhandler.CreatePermission("")
		if err != nil {
			return List{}, err
		}
		permissions_id = permission.ID
	}

	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO lists (title, description, account_id, parent_list_id, permissions_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, description, account_id, parent_list_id, permissions_id, date_created, date_edited`,
		title, description, account_id, parent_list_id, permissions_id)
	if err != nil {
		return List{}, err
	}
	var list List
	rows.Next()
	err = rows.Scan(&list.Title, &list.Description, &list.Account_ID, &list.Parent_List_ID, &list.Permissions_ID, &list.Date_Created, &list.Date_Edited)
	if err != nil {
		return List{}, err
	}
	return list, nil
}

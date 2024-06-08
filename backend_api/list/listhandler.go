package list

import (
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/todo"
	"net/http"
)

type ListHandler struct {
	db             *db.DBPool
	accounthandler *account.AccountHandler
	todohandler    *todo.TodoHandler
}

func NewListHandler(db *db.DBPool, accounthandler *account.AccountHandler, todohandler *todo.TodoHandler) *ListHandler {
	return &ListHandler{db: db, accounthandler: accounthandler, todohandler: todohandler}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {

	}
}

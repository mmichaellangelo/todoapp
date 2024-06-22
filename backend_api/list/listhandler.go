package list

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/todo"
	"net/http"
	"regexp"
	"strconv"
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
	ListRE       = regexp.MustCompile(`^\/accounts\/(\d+)\/lists\/?$`)
	ListREWithID = regexp.MustCompile(`^\/accounts\/(\d+)\/lists\/(\d+)\/?$`)
)

func NewListHandler(db *db.DBPool, accounthandler *account.AccountHandler, todohandler *todo.TodoHandler, authhandler *auth.AuthHandler) *ListHandler {
	return &ListHandler{db: db, accounthandler: accounthandler, todohandler: todohandler, authhandler: authhandler}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// Get all lists by account id
	case r.Method == http.MethodGet && ListRE.MatchString(r.URL.Path):
		groups := ListRE.FindStringSubmatch(r.URL.Path)
		if len(groups) != 2 {
			w.WriteHeader(400)
			w.Write([]byte("bad URL, error finding submatch"))
			return
		}
		id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
		}
		lists, err := h.GetAllByAccountID(id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error getting lists: " + err.Error()))
			return
		}
		resjson, err := json.Marshal(lists)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("errorrerrorerror"))
			return
		}
		w.Write(resjson)
		return
	// Get list by account ID and list ID
	case r.Method == http.MethodGet && ListREWithID.MatchString(r.URL.Path):
		fmt.Println("List by ID Route")
		groups := ListREWithID.FindStringSubmatch(r.URL.Path)
		if len(groups) != 3 {
			w.WriteHeader(500)
			w.Write([]byte("invalid request"))
		}

		account_id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("invalid account id"))
		}

		list_id, err := strconv.ParseInt(groups[2], 10, 64)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("invalid list id"))
		}
		list, err := h.GetByListID(list_id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error getting list"))
		}
		// -------------- TODO :: Check Permissions!!
		// right now only allow access to owner
		if account_id != list.Account_ID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}

		listJ, err := json.Marshal(list)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error marshalling json"))
		}
		w.Write(listJ)

	}
}

func (h *ListHandler) GetAllByAccountID(id int64) ([]ListWithTodos, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT id, title, description, account_id, parent_list_id, 
	 	permissions_id, date_created, date_edited 
	 	FROM lists
	 	WHERE account_id=$1`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lists []ListWithTodos
	for rows.Next() {
		var list ListWithTodos
		var parentListID sql.NullInt64
		var permissionsID sql.NullInt64
		err := rows.Scan(&list.ID, &list.Title, &list.Description,
			&list.Account_ID, &parentListID, &permissionsID,
			&list.Date_Created, &list.Date_Edited)
		if err != nil {
			return nil, err
		}
		if parentListID.Valid {
			list.Parent_List_ID = parentListID.Int64
		} else {
			list.Parent_List_ID = -1
		}

		if permissionsID.Valid {
			list.Permissions_ID = permissionsID.Int64
		} else {
			list.Permissions_ID = -1
		}
		lists = append(lists, list)
	}
	for i, list := range lists {
		todos, err := h.todohandler.GetAllByListID(list.ID)
		if err != nil {
			return nil, err
		}
		lists[i].Todos = todos
	}
	return lists, nil
}

func (h *ListHandler) GetByListID(id int64) (ListWithTodos, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT id, title, description, account_id, parent_list_id, 
	 	permissions_id, date_created, date_edited 
	 	FROM lists
	 	WHERE id=$1`, id)

	if err != nil {
		return ListWithTodos{}, err
	}
	defer rows.Close()
	var list ListWithTodos
	rows.Next()
	var parentListID sql.NullInt64
	var permissionsID sql.NullInt64
	err = rows.Scan(&list.ID, &list.Title, &list.Description,
		&list.Account_ID, &parentListID, &permissionsID,
		&list.Date_Created, &list.Date_Edited)
	if err != nil {
		return ListWithTodos{}, err
	}
	if parentListID.Valid {
		list.Parent_List_ID = parentListID.Int64
	} else {
		list.Parent_List_ID = -1
	}

	if permissionsID.Valid {
		list.Permissions_ID = permissionsID.Int64
	} else {
		list.Permissions_ID = -1
	}

	todos, err := h.todohandler.GetAllByListID(list.ID)
	if err != nil {
		return ListWithTodos{}, err
	}
	list.Todos = todos

	return list, nil
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

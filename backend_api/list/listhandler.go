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

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	claims, hasClaims := r.Context().Value("claims").(*auth.Claims)
	fmt.Println(claims)
	switch {
	// CREATE LIST ----------------------------------------------------
	case r.Method == http.MethodPost && ListRE.MatchString(r.URL.Path):
		// check that claims exist
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// parse url for account id
		acc_id, err := getAccountIDFromURL(r.URL.Path)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url: %v", err), http.StatusBadRequest)
			return
		}

		// query db, create list
		list, err := h.Create("", "", acc_id, -1, -1)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't create list: %v", err), http.StatusInternalServerError)
			return
		}

		// ok >> marshal data and send it back
		res, err := json.Marshal(list)
		if err != nil {
			http.Error(w, "error marshalling json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return

	// GET ALL LISTS BY ACCOUNT --------------------------------
	case r.Method == http.MethodGet && ListRE.MatchString(path):
		// before doing anything else, make sure request has claims
		if !hasClaims {
			http.Error(w, "must supply access token", http.StatusUnauthorized)
			fmt.Println("no claims")
			return
		}

		// get account id from url
		acc_id, err := getAccountIDFromURL(path)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("bad request: %v", err)))
		}

		// if account id mismatch with request, status unauthorized
		if claims.UserID != acc_id {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		// if auth ok, get account
		lists, err := h.GetAllByAccountID(acc_id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error getting lists: " + err.Error()))
			return
		}

		// write account
		resjson, err := json.Marshal(lists)
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling json: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(resjson)
		return

	// GET LIST BY ID ------------------------------------------------------
	case r.Method == http.MethodGet && ListREWithID.MatchString(r.URL.Path):
		// make sure request has claims
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// parse url, get list id
		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse url: %v", err), http.StatusBadRequest)
			return
		}

		// query db for list
		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting list: %v", err), http.StatusNotFound)
			return
		}

		// check permissions
		if claims.UserID != list.Account_ID {
			err := h.checkPermission(list.Permissions_ID, claims.UserID)
			if err != nil {
				http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
				return
			}
		}

		// all ok, marshal data and send it back
		listJ, err := json.Marshal(list)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error marshalling json"))
		}
		w.Write(listJ)
		return

	// UPDATE LIST -------------------------------------------------------------------------------------
	case (r.Method == http.MethodPatch || r.Method == http.MethodPut) && ListREWithID.MatchString(path):
		// make sure has claims
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// get id
		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url: %v", err), http.StatusBadRequest)
			return
		}

		// get list
		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting list: %v", err), http.StatusNotFound)
		}

		// check permissions
		if claims.UserID != list.Account_ID {
			err := h.checkPermission(list.Permissions_ID, claims.UserID)
			if err != nil {
				http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
				return
			}
		}

		// parse form, get new data
		err = r.ParseMultipartForm(0)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
			return
		}
		newtitle := r.FormValue("title")
		newdescription := r.FormValue("description")

		// make sure stuff has changed
		if (list.Title == newtitle) && (list.Description == newdescription) {
			// nothing changed
			w.Write([]byte("ok. nothing changed."))
			return
		}

		err = h.UpdateList(list_id, newtitle, newdescription)
		if err != nil {
			http.Error(w, fmt.Sprintf("error updating list: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	// DELETE LIST ------------------------------------------------------
	case r.Method == http.MethodDelete && ListREWithID.MatchString(path):
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't find list: %v", err), http.StatusNotFound)
		}

		if list.Account_ID != claims.UserID {
			http.Error(w, "you do not own this resource", http.StatusUnauthorized)
			return
		}

		err = h.Delete(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to delete list: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// HELPERS
func getAccountIDFromURL(url string) (int64, error) {
	groups := ListRE.FindStringSubmatch(url)
	if len(groups) != 2 {
		return -1, fmt.Errorf("invalid request")
	}

	acc_id, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid account id")
	}

	return acc_id, nil
}

func getListIDFromURL(url string) (int64, error) {
	groups := ListREWithID.FindStringSubmatch(url)
	if len(groups) != 2 {
		return -1, fmt.Errorf("invalid request")
	}

	acc_id, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid account id")
	}

	return acc_id, nil
}

func (h *ListHandler) checkPermission(perm_id int64, acc_id int64) error {
	ids, err := h.permissionhandler.GetPermissions(perm_id)
	if err != nil {
		return err
	}
	if functions.Contains(ids, acc_id) {
		return nil
	} else {
		return fmt.Errorf("unauthorized")
	}
}

// CREATE
func (h *ListHandler) Create(title string, description string, account_id int64, parent_list_id int64, permissions_id int64) (List, error) {
	// if -1 provided for permissions_id, create a new permission and use that

	if permissions_id == -1 {
		permission, err := h.permissionhandler.CreatePermission("")
		if err != nil {
			return List{}, err
		}
		permissions_id = permission.ID
	}

	var ParentListID sql.NullInt64
	if parent_list_id == -1 {
		ParentListID = sql.NullInt64{Valid: false}
	} else {
		ParentListID = sql.NullInt64{Int64: parent_list_id, Valid: true}
	}

	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO lists (title, description, account_id, parent_list_id, permissions_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, description, account_id, parent_list_id, permissions_id, date_created, date_edited`,
		title, description, account_id, ParentListID, permissions_id)
	if err != nil {
		return List{}, err
	}
	var list List
	rows.Next()
	var tempParentListID sql.NullInt64
	err = rows.Scan(&list.ID, &list.Title, &list.Description, &list.Account_ID, &tempParentListID, &list.Permissions_ID, &list.Date_Created, &list.Date_Edited)
	if err != nil {
		return List{}, err
	}
	if !tempParentListID.Valid {
		list.Parent_List_ID = -1
	} else {
		list.Parent_List_ID = tempParentListID.Int64
	}
	return list, nil
}

// READ
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

// UPDATE
func (h *ListHandler) UpdateList(id int64, newtitle string, newdesc string) error {
	_, err := h.db.Pool.Query(context.Background(),
		`UPDATE lists
		SET title=$1, description=$2
		WHERE id=$3`, newtitle, newdesc, id)
	if err != nil {
		return err
	}
	return nil
}

// DELETE
func (h *ListHandler) Delete(id int64) error {
	_, err := h.db.Pool.Query(context.Background(),
		`DELETE FROM lists
		WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

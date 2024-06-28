package todo

import (
	"context"
	"database/sql"
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
	TodoRE       = regexp.MustCompile(`^\/todos\/$`)
	TodoREWithID = regexp.MustCompile(`^\/todos\/([0-9]+)$`)
)

func NewTodoHandler(db *db.DBPool, accounthandler *account.AccountHandler) *TodoHandler {
	return &TodoHandler{db: db, accounthandler: accounthandler}
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, hasClaims := r.Context().Value("claims").(*auth.Claims)
	switch {
	// Get all
	case r.Method == http.MethodGet && TodoRE.MatchString(r.URL.Path):

	// Get by Todo ID
	case r.Method == http.MethodGet && TodoREWithID.MatchString(r.URL.Path):
		if !hasClaims {
			http.Error(w, "no access token", http.StatusUnauthorized)
			return
		}
		groups := TodoREWithID.FindStringSubmatch(r.URL.Path)
		if len(groups) != 2 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		todo_id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)
			return
		}
		t, err := h.GetByID(todo_id)
		if err != nil {
			http.Error(w, "error getting todo", http.StatusNotFound)
		}
		if !(claims.UserID == t.Account_ID) {
			//get permissions!! >> go from there
		}
	// Create
	case r.Method == http.MethodPost && TodoRE.MatchString(r.URL.Path):

	default:
		return
	}
}

// CREATE
func (h *TodoHandler) Create(body string, list_id int64, account_id int64) (Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO todos (body, list_id, completed, account_id)
			VALUES ($1, $2, $3, $4) 
			RETURNING id, body, list_id, completed, account_id, date_created, date_edited, permissions_id`,
		body, list_id, account_id)
	if err != nil {
		return Todo{}, err
	}
	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Body, &todo.List_ID,
			&todo.Completed, &todo.Account_ID, &todo.Date_Created, &todo.Date_Edited, &todo.Permissions_ID)
		if err != nil {
			return Todo{}, err
		}
	}
	return todo, nil
}

// READ
func (h *TodoHandler) GetAll() ([]Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(), `SELECT id, body, list_id, completed, 
														account_id, date_created, date_edited, permissions_id 
														FROM todos`)
	if err != nil {
		return nil, err
	}
	var todolist []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.List_ID,
			&todo.Completed, &todo.Account_ID, &todo.Date_Created, &todo.Date_Edited, &todo.Permissions_ID)
		if err != nil {
			return nil, err
		}
		todolist = append(todolist, todo)
	}
	return todolist, nil
}

func (h *TodoHandler) GetAllByAccountID(id int64) ([]Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(), `SELECT id, body, list_id, completed, 
														account_id, date_created, date_edited, permissions_id 
														FROM todos
														WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	var todolist []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.List_ID,
			&todo.Completed, &todo.Account_ID, &todo.Date_Created, &todo.Date_Edited, &todo.Permissions_ID)
		if err != nil {
			return nil, err
		}
		todolist = append(todolist, todo)
	}
	return todolist, nil
}

func (h *TodoHandler) GetByID(id int64) (Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(), `SELECT id, body, list_id, completed, 
														account_id, date_created, date_edited, permissions_id 
														FROM todos
														WHERE id=$1`, id)
	if err != nil {
		return Todo{}, err
	}

	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Body, &todo.List_ID,
			&todo.Completed, &todo.Account_ID, &todo.Date_Created, &todo.Date_Edited, &todo.Permissions_ID)
		if err != nil {
			return Todo{}, err
		}
	}
	return todo, nil
}

func (h *TodoHandler) GetAllByListID(list_id int64) ([]Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(), `SELECT id, body, list_id, completed, 
														account_id, date_created, date_edited, permissions_id 
														FROM todos
														WHERE list_id=$1`, list_id)
	if err != nil {
		return nil, err
	}
	var todolist []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Body, &todo.List_ID,
			&todo.Completed, &todo.Account_ID, &todo.Date_Created, &todo.Date_Edited, &todo.Permissions_ID)
		if err != nil {
			return nil, err
		}
		todolist = append(todolist, todo)
	}
	return todolist, nil
}

func (h *TodoHandler) GetAllByUsername(username string) ([]Todo, error) {
	account, err := h.accounthandler.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	todos, err := h.GetAllByAccountID(account.ID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (h *TodoHandler) GetAllByEmail(email string) ([]Todo, error) {
	account, err := h.accounthandler.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	todos, err := h.GetAllByAccountID(account.ID)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// UPDATE
func (h *TodoHandler) UpdateCompleted(todo_id int64, completed bool) (Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`UPDATE todos
		SET completed=$1
		WHERE id=$2
		RETURNING id, body, list_id, completed, 
		account_id, date_created, date_edited, 
		permissions_id`, completed, todo_id)
	if err != nil {
		return Todo{}, err
	}
	defer rows.Close()
	rows.Next()
	var t Todo
	err = rows.Scan(&t.ID, &t.Body, &t.List_ID,
		&t.Completed, &t.Account_ID, &t.Date_Created,
		&t.Date_Edited, &t.Permissions_ID)
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (h *TodoHandler) UpdateBody(todo_id int64, body string) (Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`UPDATE todos
		SET body=$1
		WHERE id=$2
		RETURNING id, body, list_id, completed, 
		account_id, date_created, date_edited, 
		permissions_id`, body, todo_id)
	if err != nil {
		return Todo{}, err
	}
	defer rows.Close()
	rows.Next()
	var t Todo
	err = rows.Scan(&t.ID, &t.Body, &t.List_ID,
		&t.Completed, &t.Account_ID, &t.Date_Created,
		&t.Date_Edited, &t.Permissions_ID)
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (h *TodoHandler) UpdateListID(todo_id int64, list_id int64) (Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`UPDATE todos
		SET list_id=$1
		WHERE id=$2
		RETURNING id, body, list_id, completed, 
		account_id, date_created, date_edited, 
		permissions_id`, list_id, todo_id)
	if err != nil {
		return Todo{}, err
	}
	defer rows.Close()
	rows.Next()
	var t Todo
	err = rows.Scan(&t.ID, &t.Body, &t.List_ID,
		&t.Completed, &t.Account_ID, &t.Date_Created,
		&t.Date_Edited, &t.Permissions_ID)
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

// DELETE
func (h *TodoHandler) Delete(todo_id int64) (*Todo, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`DELETE FROM todos
		WHERE id=$1
		RETURNING id, body, list_id, completed, 
		account_id, date_created, date_edited, 
		permissions_id `, todo_id)
	if err != nil {
		return nil, err
	}
	rows.Next()
	var t Todo
	err = rows.Scan(&t.ID, &t.Body, &t.List_ID,
		&t.Completed, &t.Account_ID, &t.Date_Created,
		&t.Date_Edited, &t.Permissions_ID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

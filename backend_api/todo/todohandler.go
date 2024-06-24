package todo

import (
	"context"
	"database/sql"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/db"
	"net/http"
	"regexp"
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
	TodoREWithID = regexp.MustCompile(`^\/todos\/id\/[0-9]+$`)
)

func NewTodoHandler(db *db.DBPool, accounthandler *account.AccountHandler) *TodoHandler {
	return &TodoHandler{db: db, accounthandler: accounthandler}
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// Get all
	case r.Method == http.MethodGet && TodoRE.MatchString(r.URL.Path):

	// Get by Todo ID
	case r.Method == http.MethodGet && TodoREWithID.MatchString(r.URL.Path):
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

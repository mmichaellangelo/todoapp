package todo

// HELPERS
func getAccountIDFromURL(url string) (int64, error) {
	groups := TodoRE.FindStringSubmatch(url)
	if len(groups) != 2 {
		return -1, fmt.Errorf("invalid request")
	}

	acc_id, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid account id")
	}

	return acc_id, nil
}

func getTodoIDFromURL(url string) (int64, error) {
	groups := TodoREWithID.FindStringSubmatch(url)
	if len(groups) != 2 {
		return -1, fmt.Errorf("invalid request")
	}

	acc_id, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid account id")
	}

	return acc_id, nil
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

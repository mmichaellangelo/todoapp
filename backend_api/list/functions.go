package list

import (
	"context"
	"database/sql"
	"fmt"
	"mykale/todobackendapi/functions"
	"strconv"
)

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

//
// CRUD
//

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

func (h *ListHandler) GetByID(id int64) (ListWithTodos, error) {
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
func (h *ListHandler) UpdateTitleAndDescription(id int64, newtitle string, newdesc string) error {
	_, err := h.db.Pool.Query(context.Background(),
		`UPDATE lists
		SET title=$1, description=$2
		WHERE id=$3`, newtitle, newdesc, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *ListHandler) UpdateTitle(id int64, newtitle string) error {
	_, err := h.db.Pool.Query(context.Background(),
		`UPDATE lists
		SET title=$1
		WHERE id=$2`, newtitle, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *ListHandler) UpdateDescription(id int64, newdesc string) error {
	_, err := h.db.Pool.Query(context.Background(),
		`UPDATE lists
		SET description=$1
		WHERE id=$2`, newdesc, id)
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

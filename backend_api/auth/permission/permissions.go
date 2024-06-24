package permission

import (
	"context"
	"mykale/todobackendapi/db"
)

type Permission struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PermissionHandler struct {
	db *db.DBPool
}

func NewPermissionHandler(db *db.DBPool) *PermissionHandler {
	return &PermissionHandler{db: db}
}

func (h *PermissionHandler) CreatePermission(name string) (Permission, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO permissions (name)
		VALUES ($1) RETURNING id, name`,
		name)
	if err != nil {
		return Permission{}, err
	}
	var permission Permission
	rows.Next()
	err = rows.Scan(&permission.ID, &permission.Name)
	if err != nil {
		return Permission{}, err
	}
	return permission, nil
}

func (h *PermissionHandler) EditPermissionName(id int64, newname string) (Permission, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`UPDATE permissions SET name=$1 WHERE id=$2
		 RETURNING id, name`, newname, id)
	if err != nil {
		return Permission{}, err
	}
	var permission Permission
	rows.Next()
	err = rows.Scan(&permission.ID, &permission.Name)
	if err != nil {
		return Permission{}, err
	}
	return permission, nil
}

func (h *PermissionHandler) DeletePermission(id int64) (Permission, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`DELETE FROM permissions WHERE id=$1
		 RETURNING id, name`, id)
	if err != nil {
		return Permission{}, err
	}
	var permission Permission
	rows.Next()
	err = rows.Scan(&permission.ID, &permission.Name)
	if err != nil {
		return Permission{}, err
	}
	return permission, nil
}

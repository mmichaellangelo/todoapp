package todo

import (
	"context"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/db"
	"net/http"
	"regexp"
	"time"
)

type ItemHandler struct {
	db          *db.DBPool
	authhandler *auth.AuthHandler
}

type Item struct {
	ID           int64
	ItemType     string
	Title        string
	List_ID      int64
	Completed    bool
	Account_ID   int64
	Date_Created time.Time
	Date_Edited  time.Time
}

var (
	ItemRE       = regexp.MustCompile(`^\/items\/$`)
	ItemREWithID = regexp.MustCompile(`^\/items\/id\/[0-9]+$`)
)

func NewItemHandler(db *db.DBPool, authhandler *auth.AuthHandler) http.Handler {
	return &ItemHandler{db: db, authhandler: authhandler}
}

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	// Get all items
	case r.Method == http.MethodGet && ItemRE.MatchString(r.URL.Path):

	case r.Method == http.MethodGet && ItemREWithID.MatchString(r.URL.Path):

	case r.Method == http.MethodPost && ItemRE.MatchString(r.URL.Path):

	default:
		return
	}
}

func (h *ItemHandler) CreateItem(title string, itemtype string, list_id, account_id int64) (Item, error) {
	rows, err := h.db.Pool.Query(context.Background(),
		`INSERT INTO item (title, itemtype, list_id, account_id, date_created 
			RETURNING (id, itemtype, title, list_id, account_id, date_created, date_edited)`,
		title, itemtype, list_id, account_id, time.Now())
	if err != nil {
		return Item{}, err
	}
	var item Item
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.ItemType, &item.Title, &item.List_ID,
			&item.Account_ID, &item.Date_Created, &item.Date_Edited)
		if err != nil {
			return Item{}, err
		}
	}
	item.Completed = false
	item.Date_Edited = item.Date_Created
	return item, nil
}

func (h *ItemHandler) GetAllItems() ([]Item, error) {
	rows, err := h.db.Pool.Query(context.Background(), `SELECT (id, itemtype, title, 
		list_id, completed, account_id, date_created, date_edited) FROM item`)
	if err != nil {
		return nil, err
	}
	var itemlist []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.ItemType, &item.Title, &item.List_ID,
			&item.Completed, &item.Account_ID, &item.Date_Created, &item.Date_Edited)
		if err != nil {
			return nil, err
		}
	}
	return itemlist, nil
}

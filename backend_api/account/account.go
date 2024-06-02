package account

import "time"

type Account struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"date_created"`
}

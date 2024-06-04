package account

import "time"

type Account struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHashed string    `json:"password"`
	DateCreated    time.Time `json:"date_created"`
	DateEdited     time.Time `json:"date_edited"`
}

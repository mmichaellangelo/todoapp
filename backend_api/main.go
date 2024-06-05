package main

import (
	"context"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/login"
	"mykale/todobackendapi/db"
	"net/http"
)

func main() {
	// Create DB Connection
	pool, err := db.Connect(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create handlers
	authhandler := auth.NewAuthHandler(pool)
	accounthandler := account.NewAccountHandler(pool)

	// Initialize routes
	mux := http.NewServeMux()
	mux.Handle("/accounts/", accounthandler)
	mux.Handle("/login/", login.NewLoginHandler(pool, authhandler, accounthandler))

	// Start server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error serving routes: ", err)
	}
}

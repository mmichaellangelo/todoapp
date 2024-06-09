package main

import (
	"context"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/login"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/todo"
	"net/http"
)

func main() {
	// Create DB pool
	pool, err := db.Connect(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create handlers
	authhandler := auth.NewAuthHandler(pool)
	accounthandler := account.NewAccountHandler(pool)
	todohandler := todo.NewTodoHandler(pool, authhandler, accounthandler)
	loginhandler := login.NewLoginHandler(pool, authhandler, accounthandler)

	// Initialize routes
	mux := http.NewServeMux()
	mux.Handle("/accounts/", accounthandler)
	mux.Handle("/login/", loginhandler)
	mux.Handle("/todos/", todohandler)

	// Start server
	err = http.ListenAndServe(":80", mux)
	if err != nil {
		fmt.Println("Error serving routes: ", err)
	}
}

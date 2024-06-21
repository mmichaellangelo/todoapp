package main

import (
	"context"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/login"
	"mykale/todobackendapi/combined"
	"mykale/todobackendapi/db"
	"mykale/todobackendapi/list"
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

	// Main handlers
	authhandler := auth.NewAuthHandler(pool)
	accounthandler := account.NewAccountHandler(pool, authhandler)
	todohandler := todo.NewTodoHandler(pool, authhandler, accounthandler)
	listhandler := list.NewListHandler(pool, accounthandler, todohandler, authhandler)
	loginhandler := login.NewLoginHandler(pool, authhandler, accounthandler)

	// Combined handler delegates routes from /account/* to respective handlers
	combinedhandler := combined.NewCombinedHandler(accounthandler, listhandler, todohandler)

	// Initialize routes
	mux := http.NewServeMux()
	mux.Handle("/login/", loginhandler)
	mux.Handle("/accounts/", combinedhandler)

	// ------------- TODO: Logging Middleware

	// Start server
	err = http.ListenAndServe(":80", mux)
	if err != nil {
		fmt.Println("Error serving routes: ", err)
	}
}

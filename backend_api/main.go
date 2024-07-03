package main

import (
	"context"
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/auth"
	"mykale/todobackendapi/auth/login"
	"mykale/todobackendapi/auth/permission"
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
	permissionhandler := permission.NewPermissionHandler(pool)
	accounthandler := account.NewAccountHandler(pool)
	todohandler := todo.NewTodoHandler(pool, accounthandler)
	listhandler := list.NewListHandler(pool, accounthandler, todohandler, permissionhandler)
	loginhandler := login.NewLoginHandler(pool, authhandler, accounthandler)

	// Combined handler delegates routes from /account/* to respective handlers
	combinedhandler := combined.NewCombinedHandler(accounthandler, listhandler, todohandler)

	// Initialize routes
	mux := http.NewServeMux()
	mux.Handle("/login/", loginhandler)
	mux.Handle("/logout/", loginhandler)
	// combined handler delegates acccounts/, accounts/[id]/{resource} to respective handlers
	mux.Handle("/accounts/", combinedhandler)
	mux.Handle("/todos/", todohandler)
	mux.Handle("/lists/", listhandler)

	// Logger Middleware
	LoggerMux := NewLoggerMiddleware(mux)
	defer LoggerMux.Close()

	// Auth Middleware
	AuthMux := auth.NewAuthMiddleware(LoggerMux)

	// Start server
	err = http.ListenAndServe(":80", AuthMux)
	if err != nil {
		fmt.Println("Error serving routes: ", err)
	}

}

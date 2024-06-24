package auth

import (
	"fmt"
	"mykale/todobackendapi/account"
	"mykale/todobackendapi/list"
	"mykale/todobackendapi/todo"
	"net/http"
	"regexp"
)

type AuthMiddleware struct {
	accounthandler *account.AccountHandler
	listhandler    *list.ListHandler
	todohandler    *todo.TodoHandler
	next           http.Handler
}

func NewAuthMiddleware(handlerToWrap http.Handler) *AuthMiddleware {
	return &AuthMiddleware{next: handlerToWrap}
}

func (h *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isRestrictedPath(r.URL.Path) {
		fmt.Println("RESTRICTED PATH")
		accesstoken := r.Header.Get("accesstoken")
		cookies := r.Cookies()
		if cookies == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, err := GetClaimsFromToken(accesstoken)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("CLAIMS: Username: %v, UserID: %d\n", claims.Username, claims.UserID)
	}
	h.next.ServeHTTP(w, r)
	fmt.Println("auth middleware")
}

func isRestrictedPath(path string) bool {
	PathRE := regexp.MustCompile(`^\/accounts\/.*$`)
	if PathRE.MatchString(path) {
		return true
	} else {
		return false
	}
}

func isLoginPath(path string) bool {
	LoginRE := regexp.MustCompile(`^\/login\/?$`)
	if LoginRE.MatchString(path) {
		return true
	} else {
		return false
	}
}

// Access Protected Resource
// -- get user info from access token
// -- get permission info from resource
// -- cross-check permissions with user info
// -- grant access or throw error

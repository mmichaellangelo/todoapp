package auth

import (
	"fmt"
	"net/http"
	"regexp"
)

type AuthMiddleware struct {
	handler http.Handler
}

func NewAuthMiddleware(handlerToWrap http.Handler) *AuthMiddleware {
	return &AuthMiddleware{handler: handlerToWrap}
}

func (h *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isRestrictedPath(r.URL.Path) {
		fmt.Println("RESTRICTED PATH")
	}
	h.handler.ServeHTTP(w, r)
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

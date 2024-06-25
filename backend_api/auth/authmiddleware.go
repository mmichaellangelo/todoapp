package auth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

type AuthMiddleware struct {
	next http.Handler
}

func NewAuthMiddleware(handlerToWrap http.Handler) *AuthMiddleware {
	return &AuthMiddleware{next: handlerToWrap}
}

func (h *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isRestrictedPath(r.URL.Path) {
		fmt.Println("RESTRICTED PATH")
		accesstoken := r.Header.Get("accesstoken")
		if accesstoken == "" {
			http.Error(w, "missing access token", http.StatusUnauthorized)
			return
		}

		claims, err := GetClaimsFromToken(accesstoken)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid token: %v", err), http.StatusBadRequest)
			return
		}

		fmt.Printf("CLAIMS: Username: %v, UserID: %d\n", claims.Username, claims.UserID)

		// add claims to context
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)
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

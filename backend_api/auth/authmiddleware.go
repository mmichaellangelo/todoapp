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

var (
	RestrictedPathRE = regexp.MustCompile(`^\/accounts\/.*$`)
	LoginPathRE      = regexp.MustCompile(`^\/login\/?$`)
)

func (h *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	// Restricted Path
	case RestrictedPathRE.MatchString(path):
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
}

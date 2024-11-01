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
		fmt.Println(accesstoken)
		if accesstoken == "undefined" || accesstoken == "" {
			fmt.Println("no access token provided")
		} else {
			claims, err := GetClaimsFromToken(accesstoken)
			if err != nil {
				fmt.Printf("error getting claims: %v", err)
			}

			fmt.Printf("CLAIMS: Username: %v, UserID: %d\n", claims.Username, claims.UserID)
			// add claims to context
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)
		}
	}
	h.next.ServeHTTP(w, r)
}

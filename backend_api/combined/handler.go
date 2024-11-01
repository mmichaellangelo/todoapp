package combined

import (
	"net/http"
	"regexp"
)

type CombinedHandler struct {
	accounthandler http.Handler
	listhandler    http.Handler
	todohandler    http.Handler
}

var (
	ListRE    = regexp.MustCompile(`^\/accounts\/\d+\/lists\/?$`)
	TodoRE    = regexp.MustCompile(`^\/accounts\/\d+\/todos\/?$`)
	AccountRE = regexp.MustCompile(`^\/accounts\/?(\d+\/?)?$`)
)

// NewCombinedHandler creates a new CombinedHandler.
func NewCombinedHandler(accounthandler, listhandler, todohandler http.Handler) *CombinedHandler {
	return &CombinedHandler{
		accounthandler: accounthandler,
		listhandler:    listhandler,
		todohandler:    todohandler,
	}
}

// ServeHTTP routes the requests to the appropriate sub-handler.
func (h *CombinedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	// List Route
	case ListRE.MatchString(path):
		h.listhandler.ServeHTTP(w, r)
		return
	// Todo Route
	case TodoRE.MatchString(path):
		h.todohandler.ServeHTTP(w, r)
		return
	// Account Route
	case AccountRE.MatchString(path):
		h.accounthandler.ServeHTTP(w, r)
		return
	}

}

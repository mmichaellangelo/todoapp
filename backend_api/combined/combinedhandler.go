package combined

import (
	"fmt"
	"net/http"
	"regexp"
)

type CombinedHandler struct {
	accounthandler http.Handler
	listhandler    http.Handler
	todohandler    http.Handler
}

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
	fmt.Println("combined handler")
	if match, _ := regexp.MatchString(`^\/accounts\/\d+\/lists\/?$`, path); match {
		fmt.Println("Lists!")
		h.listhandler.ServeHTTP(w, r)
	} else if match, _ := regexp.MatchString(`^/accounts/\d+/lists/\d+/todos`, path); match {
		h.todohandler.ServeHTTP(w, r)
	} else {
		h.accounthandler.ServeHTTP(w, r)
	}
}

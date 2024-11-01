package todo

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get claims
	claims, hasClaims := r.Context().Value("claims").(*auth.Claims)
	switch {
	// ------------------------------TODO: Need routes for get all todos by account AND all todos by list
	// Get all
	case r.Method == http.MethodGet && TodoRE.MatchString(r.URL.Path):

	// Get by Todo ID
	case r.Method == http.MethodGet && TodoREWithID.MatchString(r.URL.Path):
		if !hasClaims {
			http.Error(w, "no access token", http.StatusUnauthorized)
			return
		}
		groups := TodoREWithID.FindStringSubmatch(r.URL.Path)
		if len(groups) != 2 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		todo_id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid id format", http.StatusBadRequest)
			return
		}
		t, err := h.GetByID(todo_id)
		if err != nil {
			http.Error(w, "error getting todo", http.StatusNotFound)
		}
		if !(claims.UserID == t.Account_ID) {
			//get permissions!! >> go from there
		}
	// Create
	case r.Method == http.MethodPost && TodoRE.MatchString(r.URL.Path):

	default:
		return
	}
}
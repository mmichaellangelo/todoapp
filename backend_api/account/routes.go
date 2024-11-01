package account

// ROUTES
func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	claims, hasClaims := r.Context().Value("claims").(*auth.Claims)
	if hasClaims {
		fmt.Printf("Claims: %v\n", claims)
	}
	switch {
	//----------------------------------- GET ALL ACCOUNTS
	case r.Method == http.MethodGet && AccountRE.MatchString(r.URL.Path):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// MUST HAVE API KEY ----------- ** TODO
		accounts, err := h.GetAll()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		resp, err := json.Marshal(accounts)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(resp)
		return
		//----------------------------------- GET ACCOUNT BY ID
	case r.Method == http.MethodGet && AccountREWithID.MatchString(r.URL.Path):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		groups := AccountREWithID.FindStringSubmatch(r.URL.Path)
		if len(groups) != 2 {
			w.WriteHeader(400)
			w.Write([]byte("invalid request"))
			return
		}
		id, err := strconv.ParseInt(groups[1], 10, 64)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error parsing integer"))
			return
		}
		account, err := h.GetByID(id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		resp, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(resp)
		return
	// ----------------------------- CREATE ACCOUNT
	case r.Method == http.MethodPost && AccountRE.MatchString(r.URL.Path):
		err := r.ParseMultipartForm(0)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		validemail, err := mail.ParseAddress(email)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid email address"))
		}

		account, err := h.Create(username, validemail.Address, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(account)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
		return

	default:
		return
	}
}
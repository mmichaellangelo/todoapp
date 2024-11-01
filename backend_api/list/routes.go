package list

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	claims, hasClaims := r.Context().Value("claims").(*auth.Claims)
	fmt.Println(claims)
	switch {
	// CREATE LIST ----------------------------------------------------
	case r.Method == http.MethodPost && ListRE.MatchString(r.URL.Path):
		// check that claims exist
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// parse url for account id
		acc_id, err := getAccountIDFromURL(r.URL.Path)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url: %v", err), http.StatusBadRequest)
			return
		}

		// query db, create list
		list, err := h.Create("", "", acc_id, -1, -1)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't create list: %v", err), http.StatusInternalServerError)
			return
		}

		// ok >> marshal data and send it back
		res, err := json.Marshal(list)
		if err != nil {
			http.Error(w, "error marshalling json", http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return

	// GET ALL LISTS BY ACCOUNT --------------------------------
	case r.Method == http.MethodGet && ListRE.MatchString(path):
		// before doing anything else, make sure request has claims
		if !hasClaims {
			http.Error(w, "must supply access token", http.StatusUnauthorized)
			fmt.Println("no claims")
			return
		}

		// get account id from url
		acc_id, err := getAccountIDFromURL(path)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("bad request: %v", err)))
		}

		// if account id mismatch with request, status unauthorized
		if claims.UserID != acc_id {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		// if auth ok, get account
		lists, err := h.GetAllByAccountID(acc_id)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error getting lists: " + err.Error()))
			return
		}

		// write account
		resjson, err := json.Marshal(lists)
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling json: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(resjson)
		return

	// GET LIST BY ID ------------------------------------------------------
	case r.Method == http.MethodGet && ListREWithID.MatchString(r.URL.Path):
		// make sure request has claims
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// parse url, get list id
		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse url: %v", err), http.StatusBadRequest)
			return
		}

		// query db for list
		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting list: %v", err), http.StatusNotFound)
			return
		}

		// check permissions
		if claims.UserID != list.Account_ID {
			err := h.checkPermission(list.Permissions_ID, claims.UserID)
			if err != nil {
				http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
				return
			}
		}

		// all ok, marshal data and send it back
		listJ, err := json.Marshal(list)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error marshalling json"))
		}
		w.Write(listJ)
		return

	// UPDATE LIST -------------------------------------------------------------------------------------
	case (r.Method == http.MethodPatch || r.Method == http.MethodPut) && ListREWithID.MatchString(path):
		// make sure has claims
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		// get list id
		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url: %v", err), http.StatusBadRequest)
			return
		}

		// get list
		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting list: %v", err), http.StatusNotFound)
		}

		// check permissions
		if claims.UserID != list.Account_ID {
			err := h.checkPermission(list.Permissions_ID, claims.UserID)
			if err != nil {
				http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
				return
			}
		}

		// parse form, get new data
		err = r.ParseMultipartForm(0)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
			return
		}
		newtitle := r.FormValue("title")
		newdescription := r.FormValue("description")

		// make sure stuff has changed
		if (list.Title == newtitle) && (list.Description == newdescription) {
			// nothing changed
			w.Write([]byte("ok. nothing changed."))
			return
		}

		err = h.UpdateList(list_id, newtitle, newdescription)
		if err != nil {
			http.Error(w, fmt.Sprintf("error updating list: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	// DELETE LIST ------------------------------------------------------
	case r.Method == http.MethodDelete && ListREWithID.MatchString(path):
		if !hasClaims {
			http.Error(w, "must provide credentials", http.StatusUnauthorized)
			return
		}

		list_id, err := getListIDFromURL(path)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		list, err := h.GetByListID(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't find list: %v", err), http.StatusNotFound)
		}

		if list.Account_ID != claims.UserID {
			http.Error(w, "you do not own this resource", http.StatusUnauthorized)
			return
		}

		err = h.Delete(list_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to delete list: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}
}
package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"

	"golang.org/x/crypto/bcrypt"
)

func loginHandlerFunc(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		req.ParseForm()
		user := &api.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}

		gs := gate.NewGateService("", "")
		dbUser, err := gs.GetUser(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := newAuthToken(dbUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenCookie := newAccessCookie(token)
		http.SetCookie(w, tokenCookie)
		// TODO FIND REDIRECT AND MAKE ERRORS ACTUALLY SHOW SOMETHIN
		w.WriteHeader(http.StatusFound)
	}

	if req.Method == http.MethodGet {
		_, err := getUserFromCookie(req)

		if err == nil {
			http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
			return
		}

		err = tpl.ExecuteTemplate(w, "login", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

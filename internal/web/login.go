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

		token, err := newSession(userToProfile(dbUser))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenCookie := newAccessCookie(token.Id)
		http.SetCookie(w, tokenCookie)
		http.Redirect(w, req, req.FormValue("referer"), http.StatusFound)
	}

	if req.Method == http.MethodGet {
		_, err := getProfFromCookie(req)

		referer := req.Referer()

		if err == nil {
			http.Redirect(w, req, referer, http.StatusFound)
			return
		}

		err = tpl.ExecuteTemplate(w, "login", referer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

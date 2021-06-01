package web

import (
	"net/http"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/bennycio/bundle/logger"
)

func signupHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		user := &api.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		gs := gate.NewGateService("", "")

		err := gs.InsertUser(user)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		dbUser, err := gs.GetUser(user)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		token, err := newSession(userToProfile(dbUser))
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		tokenCookie := newAccessCookie(token.Id)
		http.SetCookie(w, tokenCookie)

		ref := r.FormValue("referer")
		if strings.Contains(ref, "login") || strings.Contains(ref, "signup") {
			ref = "/"
		}
		http.Redirect(w, r, ref, http.StatusFound)
		return
	}

	referer := r.Referer()

	err := tpl.ExecuteTemplate(w, "register", templateData{Referrer: referer})
	if err != nil {
		logger.ErrLog.Print(err.Error())
	}

}

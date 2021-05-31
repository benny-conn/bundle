package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/bennycio/bundle/logger"

	"golang.org/x/crypto/bcrypt"
)

func loginHandlerFunc(w http.ResponseWriter, req *http.Request) {
	referer := req.Referer()

	if req.Method == http.MethodPost {

		req.ParseForm()
		user := &api.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}

		gs := gate.NewGateService("", "")
		dbUser, err := gs.GetUser(user)

		if err != nil {
			err = tpl.ExecuteTemplate(w, "login", templateData{Referrer: referer, Error: errorData{
				Code:    http.StatusUnauthorized,
				Message: cleanError(err).Error(),
			}})
			if err != nil {
				logger.ErrLog.Panic(err.Error())
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

		if err != nil {
			err = tpl.ExecuteTemplate(w, "login", templateData{Referrer: referer, Error: errorData{
				Code:    http.StatusUnauthorized,
				Message: cleanError(err).Error(),
			}})
			if err != nil {
				logger.ErrLog.Panic(err.Error())
			}
			return
		}

		token, err := newSession(userToProfile(dbUser))
		if err != nil {
			err = tpl.ExecuteTemplate(w, "login", templateData{Referrer: referer, Error: errorData{
				Code:    http.StatusUnauthorized,
				Message: cleanError(err).Error(),
			}})
			if err != nil {
				logger.ErrLog.Panic(err.Error())
			}
			return
		}
		tokenCookie := newAccessCookie(token.Id)
		http.SetCookie(w, tokenCookie)
		http.Redirect(w, req, req.FormValue("referer"), http.StatusFound)
	}

	if req.Method == http.MethodGet {
		_, err := getProfFromCookie(req)

		if err == nil {
			http.Redirect(w, req, referer, http.StatusFound)
			return
		}

		err = tpl.ExecuteTemplate(w, "login", templateData{Referrer: referer})
		if err != nil {
			logger.ErrLog.Panic(err.Error())
		}
	}
}

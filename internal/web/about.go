package web

import (
	"net/http"

	"github.com/bennycio/bundle/logger"
)

func aboutHandlerFunc(w http.ResponseWriter, r *http.Request) {

	prof, err := getProfFromCookie(r)
	data := templateData{}
	if err == nil {
		data.Profile = prof
	}

	err = tpl.ExecuteTemplate(w, "about", data)
	if err != nil {
		logger.ErrLog.Print(err.Error())
	}
}

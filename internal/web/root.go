package web

import (
	"net/http"

	"github.com/bennycio/bundle/internal/logger"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	prof, err := getProfFromCookie(r)
	data := templateData{}
	if err == nil {
		data.Profile = prof
	} else {
		c, err := r.Cookie("access_token")
		if err == nil {
			c.MaxAge = -1
			http.SetCookie(w, c)
		}
	}

	err = tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		logger.ErrLog.Panic(err.Error())
	}
}

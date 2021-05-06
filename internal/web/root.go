package web

import (
	"net/http"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	prof, err := getProfFromCookie(r)
	data := TemplateData{}
	if err == nil {
		data.Profile = prof
	}

	err = tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

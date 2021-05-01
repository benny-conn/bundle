package web

import (
	"net/http"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromCookie(r)
	data := TemplateData{}
	if err == nil {
		data.User = user
	}

	err = tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

package web

import (
	"net/http"
)

func rootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	user, _ := getUserFromCookie(r)
	data := TemplateData{
		User: user,
	}

	err := tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

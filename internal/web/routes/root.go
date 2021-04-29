package routes

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
)

func RootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	user, _ := getProfileFromCookie(r)
	data := bundle.TemplateData{
		User: *user,
	}

	err := tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

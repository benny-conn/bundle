package routes

import (
	"net/http"
)

func RootHandlerFunc(w http.ResponseWriter, r *http.Request) {

	user, _ := getProfileFromCookie(r)
	data := TemplateData{
		Profile: user,
	}

	err := tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}
}

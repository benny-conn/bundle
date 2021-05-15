package web

import (
	"net/http"
)

func profileHandlerFunc(w http.ResponseWriter, req *http.Request) {

	pro, err := getProfFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	data := templateData{
		Profile: pro,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

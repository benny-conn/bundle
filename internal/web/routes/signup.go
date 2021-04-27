package routes

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
)

func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, _ := getProfileFromCookie(r)
	err := tpl.ExecuteTemplate(w, "register", user)
	if err != nil {
		bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

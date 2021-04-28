package routes

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/wrapper"
)

func UsersHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		newUser := &api.User{
			Username: req.FormValue("username"),
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		}
		err := wrapper.InsertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

}

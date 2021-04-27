package routes

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
)

func UsersHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		newUser := bundle.User{
			Username: req.FormValue("username"),
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		}
		err := storage.InsertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

}

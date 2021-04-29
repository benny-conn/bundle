package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
)

func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()
		port := os.Getenv("API_PORT")
		u, err := url.Parse(fmt.Sprintf(":%v/api/users", port))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := &api.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		userJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		buf.Write(userJSON)

		resp, err := http.Post(u.String(), "application/json", buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	user, _ := getProfileFromCookie(r)
	err := tpl.ExecuteTemplate(w, "register", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

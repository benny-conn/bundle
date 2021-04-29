package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
)

func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()
		port := os.Getenv("API_PORT")
		u, err := url.Parse(fmt.Sprintf("http://localhost:%v/api/users", port))
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
		bs, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bs))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	user, _ := getProfileFromCookie(r)

	td := internal.TemplateData{
		Profile: user,
	}

	err := tpl.ExecuteTemplate(w, "register", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

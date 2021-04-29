package routes

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	auth "github.com/bennycio/bundle/internal/auth/user"
	"github.com/bennycio/bundle/wrapper"
)

func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		r.ParseForm()

		user := &api.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		err := wrapper.InsertUserApi(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := auth.NewAuthToken(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenCookie := auth.NewAccessCookie(token)
		http.SetCookie(w, tokenCookie)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	user, _ := getProfileFromCookie(r)

	td := TemplateData{
		Profile: user,
	}

	err := tpl.ExecuteTemplate(w, "register", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	auth "github.com/bennycio/bundle/internal/auth/user"
	"github.com/bennycio/bundle/wrapper"
	"golang.org/x/crypto/bcrypt"
)

func loginHandlerFunc(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		req.ParseForm()
		user := &api.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}

		isValid := internal.IsUserValid(user)

		if !isValid {
			http.Error(w, "invalid request format", http.StatusBadRequest)
			return
		}

		dbUser, err := wrapper.GetUserApi(user.Username, user.Email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := auth.NewAuthToken(dbUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenCookie := auth.NewAccessCookie(token)
		http.SetCookie(w, tokenCookie)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	}

	if req.Method == http.MethodGet {
		user, _ := getProfileFromCookie(req)

		td := TemplateData{
			Profile: user,
		}

		err := tpl.ExecuteTemplate(w, "login", td)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

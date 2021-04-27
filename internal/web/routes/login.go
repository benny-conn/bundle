package routes

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		user := bundle.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}

		isValid := bundle.IsUserValid(user)

		if !isValid {
			http.Error(w, "invalid request format", http.StatusBadRequest)
			return
		}

		dbUser, err := storage.GetUser(user)

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
		err := tpl.ExecuteTemplate(w, "login", user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

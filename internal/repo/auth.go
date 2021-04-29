package repo

import (
	"encoding/json"
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
	"golang.org/x/crypto/bcrypt"
)

func authUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			userJSON := r.Header.Get("User")
			user := &api.User{}

			err := json.Unmarshal([]byte(userJSON), user)
			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			isValid := internal.IsUserValid(user)

			if !isValid {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			dbUser, err := wrapper.GetUserApi(user.Username, user.Email)

			if err != nil {
				http.Error(w, "invalid user", http.StatusBadRequest)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))

			if err != nil {
				http.Error(w, "incorrect password", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

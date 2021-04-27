package routes

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
)

func getProfileFromCookie(r *http.Request) (bundle.User, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return bundle.User{}, err
	}
	user, err := auth.GetUserFromToken(c.Value)
	if err != nil {
		return bundle.User{}, err
	}
	return user, nil

}

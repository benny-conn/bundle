package routes

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/auth"
)

func getProfileFromCookie(r *http.Request) (*api.User, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	user, err := auth.GetUserFromToken(c.Value)
	if err != nil {
		return nil, err
	}
	return user, nil

}

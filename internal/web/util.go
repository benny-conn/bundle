package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
)

func getUserFromCookie(r *http.Request) (*api.User, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return nil, err
	}
	user, err := getUserFromToken(c.Value)
	if err != nil {
		return nil, err
	}
	return user, nil

}

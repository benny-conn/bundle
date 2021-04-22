package web

import (
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
)

func getProfileFromCookie(r *http.Request) (bundle.Profile, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return bundle.Profile{}, err
	}
	user, err := auth.GetProfileFromToken(c.Value)
	if err != nil {
		return bundle.Profile{}, err
	}
	return user, nil

}

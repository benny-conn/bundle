package web

import (
	"net/http"
)

func getProfFromCookie(r *http.Request) (profile, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return profile{}, err
	}
	user, err := getProfileFromToken(c.Value)
	if err != nil {
		return profile{}, err
	}
	return user, nil

}

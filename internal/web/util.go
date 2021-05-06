package web

import (
	"net/http"
)

func getProfFromCookie(r *http.Request) (Profile, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return Profile{}, err
	}
	user, err := getProfileFromToken(c.Value)
	if err != nil {
		return Profile{}, err
	}
	return user, nil

}

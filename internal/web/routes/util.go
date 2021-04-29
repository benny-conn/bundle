package routes

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
)

func getProfileFromCookie(r *http.Request) (internal.Profile, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return internal.Profile{}, err
	}
	user, err := auth.GetUserFromToken(c.Value)
	if err != nil {
		return internal.Profile{}, err
	}
	return userToProfile(user), nil

}

func userToProfile(user *api.User) internal.Profile {
	return internal.Profile{
		Username: user.Username,
		Email:    user.Email,
		Tag:      user.Tag,
	}
}

package orm

import (
	"regexp"

	"github.com/bennycio/bundle/api"
)

func isUserValid(user *api.User) bool {

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.Email != "" {
		if len(user.Email) > 254 || !rxEmail.MatchString(user.Email) {
			return false
		}
	}

	return true

}

package internal

import (
	"os"
	"regexp"
)

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsUserValid(user User) bool {

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.Email != "" {
		if len(user.Email) > 254 || !rxEmail.MatchString(user.Email) {
			return false
		}
	}

	return true

}

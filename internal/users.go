package internal

import (
	"encoding/json"
	"errors"
	"regexp"
)

func ValidateAndReturnUser(userAsJSON string) (*User, error) {

	if userAsJSON == "" || !IsJSON(userAsJSON) {
		return nil, errors.New("must be in JSON format")
	}

	u := &User{}

	err := json.Unmarshal([]byte(userAsJSON), u)

	if err != nil {
		return nil, err
	}

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if u.Email != "" {
		if len(u.Email) > 254 || !rxEmail.MatchString(u.Email) {
			return nil, errors.New("invalid email")
		}
	}

	return u, nil

}

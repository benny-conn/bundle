package internal

import (
	"encoding/json"
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
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

func GetUserFromDatabase(user *User) (*User, error) {
	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	decodedUser := &User{}

	if user.Email == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"username", user.Username}}).Decode(decodedUser)
	} else if user.Username == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"email", NewCaseInsensitiveRegex(user.Email)}}).Decode(decodedUser)
	} else {
		err = collection.FindOne(session.Ctx, bson.D{{"username", user.Username}, {"email", NewCaseInsensitiveRegex(user.Email)}}).Decode(decodedUser)
	}
	if err != nil {
		return nil, err
	}

	return decodedUser, nil
}

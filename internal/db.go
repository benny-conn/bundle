package internal

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user User) error {
	asJSON, _ := json.Marshal(user)
	validatedUser, err := ValidateAndReturnUser(string(asJSON))

	if err != nil {
		return err
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(validatedUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	session, err := GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", NewCaseInsensitiveRegex(validatedUser.Username)}})

	if err != nil {
		return err
	}

	if countUserName > 0 {
		err = errors.New("user already exists with given username")
		return err
	}

	countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", NewCaseInsensitiveRegex(validatedUser.Email)}})

	if err != nil {
		return err
	}

	if countEmail > 0 {
		err = errors.New("user already exists with given email")
		return err
	}

	_, err = collection.InsertOne(session.Ctx, bson.D{{"username", validatedUser.Username}, {"email", validatedUser.Email}, {"password", string(bcryptPass)}})

	if err != nil {
		return err
	}
	return nil
}

func GetUser(user User) (*User, error) {
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

func InsertPlugin(plugin Plugin) error {
	session, err := GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	_, err = collection.InsertOne(session.Ctx, bson.D{{"plugin", plugin.Plugin}, {"user", plugin.User}, {"version", plugin.Version}})

	if err != nil {
		return err
	}
	return nil

}

func UpdatePlugin(name string, plugin Plugin) error {

	session, err := GetMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(session.Ctx, bson.D{{"plugin", NewCaseInsensitiveRegex(name)}}, bson.D{{"$plugin", decodedPluginResult.Plugin}, {"$user", plugin.User}, {"$version", plugin.Version}})
	if err != nil {
		return err
	}
	return nil

}

func GetPlugin(name string) (*Plugin, error) {

	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return nil, err
	}

	return decodedPluginResult, nil

}

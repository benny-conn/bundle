package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user *api.User) error {
	isValid := bundle.IsUserValid(user)

	if !isValid {
		return errors.New("invalid user")
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(bcryptPass)

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", bundle.NewCaseInsensitiveRegex(user.Username)}})

	if err != nil {
		return err
	}

	if countUserName > 0 {
		err = errors.New("user already exists with given username")
		return err
	}

	countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", bundle.NewCaseInsensitiveRegex(user.Email)}})

	if err != nil {
		return err
	}

	if countEmail > 0 {
		err = errors.New("user already exists with given email")
		return err
	}

	userToInsert := marshallBsonClean(user)

	_, err = collection.InsertOne(session.Ctx, userToInsert)

	if err != nil {
		return err
	}
	return nil
}

func GetUser(username string, email string) (*api.User, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	decodedUser := &api.User{}

	if email == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"username", username}}).Decode(decodedUser)
	} else if username == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"email", bundle.NewCaseInsensitiveRegex(email)}}).Decode(decodedUser)
	} else {
		err = collection.FindOne(session.Ctx, bson.D{{"username", username}, {"email", bundle.NewCaseInsensitiveRegex(email)}}).Decode(decodedUser)
	}
	if err != nil {
		return nil, err
	}

	return decodedUser, nil
}

func UpdateUser(username string, user *api.User) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	updatedUser := marshallBsonClean(user)

	updateResult, err := collection.UpdateOne(session.Ctx, bson.D{{"username", bundle.NewCaseInsensitiveRegex(username)}}, bson.D{{"$set", updatedUser}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no user found")
	}
	return nil
}

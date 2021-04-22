package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(user bundle.User) error {
	asJSON, _ := json.Marshal(user)
	validatedUser, err := bundle.ValidateAndReturnUser(string(asJSON))

	if err != nil {
		return err
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(validatedUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", bundle.NewCaseInsensitiveRegex(validatedUser.Username)}})

	if err != nil {
		return err
	}

	if countUserName > 0 {
		err = errors.New("user already exists with given username")
		return err
	}

	countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", bundle.NewCaseInsensitiveRegex(validatedUser.Email)}})

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

func GetUser(user bundle.User) (bundle.User, error) {
	session, err := getMongoSession()
	if err != nil {
		return bundle.User{}, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	decodedUser := &bundle.User{}

	if user.Email == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"username", user.Username}}).Decode(decodedUser)
	} else if user.Username == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"email", bundle.NewCaseInsensitiveRegex(user.Email)}}).Decode(decodedUser)
	} else {
		err = collection.FindOne(session.Ctx, bson.D{{"username", user.Username}, {"email", bundle.NewCaseInsensitiveRegex(user.Email)}}).Decode(decodedUser)
	}
	if err != nil {
		return bundle.User{}, err
	}

	return *decodedUser, nil
}

func InsertPlugin(plugin bundle.Plugin) error {
	session, err := getMongoSession()
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

func UpdatePlugin(name string, plugin bundle.Plugin) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &bundle.Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", bundle.NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(session.Ctx, bson.D{{"plugin", bundle.NewCaseInsensitiveRegex(name)}}, bson.D{{"$plugin", decodedPluginResult.Plugin}, {"$user", plugin.User}, {"$version", plugin.Version}})
	if err != nil {
		return err
	}
	return nil

}

func GetPlugin(name string) (bundle.Plugin, error) {

	session, err := getMongoSession()
	if err != nil {
		return bundle.Plugin{}, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &bundle.Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", bundle.NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return bundle.Plugin{}, err
	}

	return *decodedPluginResult, nil

}

func getMongoSession() (*bundle.Mongo, error) {
	mg := &bundle.Mongo{}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MongoURL")))
	mg.Cancel = cancel
	mg.Client = client
	mg.Ctx = ctx
	if err != nil {
		return mg, err
	}
	return mg, nil
}

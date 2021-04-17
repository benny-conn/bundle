package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ValidateAndReturnUser(user string) (*User, error) {

	if user == "" || !IsJSON(user) {
		return nil, errors.New("must be in JSON format")
	}

	u := &User{}

	err := json.Unmarshal([]byte(user), u)

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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://benny-bundle:thisismypassword1@bundle.mveuj.mongodb.net/users?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("users").Collection("users")

	var dbUser *mongo.SingleResult

	if user.Email == "" {
		dbUser = collection.FindOne(ctx, bson.D{{"username", user.Username}})
	} else if user.Username == "" {
		dbUser = collection.FindOne(ctx, bson.D{{"email", user.Email}})
	} else {
		dbUser = collection.FindOne(ctx, bson.D{{"username", user.Username}, {"email", user.Email}})
	}

	if dbUser.Err() != nil {
		return nil, dbUser.Err()
	}

	decodedUser := &User{}
	err = dbUser.Decode(decodedUser)
	if err != nil {
		return nil, err
	}
	return decodedUser, nil
}

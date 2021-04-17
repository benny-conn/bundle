package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPluginAuthor(pluginName string) (string, error) {

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

	collection := client.Database("users").Collection("plugins")

	var pluginResult *mongo.SingleResult

	pluginResult = collection.FindOne(ctx, bson.D{{"plugin", pluginName}})

	if pluginResult.Err() != nil {
		return "", pluginResult.Err()
	}

	decodedPluginResult := &struct {
		Plugin string `json:"plugin"`
		User   string `json:"user"`
	}{}
	err = pluginResult.Decode(decodedPluginResult)
	if err != nil {
		return "", err
	}
	return decodedPluginResult.User, nil
}

func GetPluginVersion(pluginName string) (string, error) {
	return "1.0.0", nil
}

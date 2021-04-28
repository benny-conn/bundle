package orm

import (
	"errors"
	"time"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertPlugin(plugin *api.Plugin) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")

	newPlugin := marshallBsonClean(plugin)
	newPlugin = append(newPlugin, bson.E{"lastUpdated", time.Now().Unix()})

	_, err = collection.InsertOne(session.Ctx, newPlugin)

	if err != nil {
		return err
	}
	return nil

}

func UpdatePlugin(name string, plugin *api.Plugin) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")

	updatedPlugin := marshallBsonClean(plugin)
	updatedPlugin = append(updatedPlugin, bson.E{"lastUpdated", time.Now().Unix()})

	updateResult, err := collection.UpdateOne(session.Ctx, bson.D{{"name", bundle.NewCaseInsensitiveRegex(name)}}, bson.D{{"$set", updatedPlugin}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no plugin found")
	}
	return nil

}

func GetPlugin(name string) (*api.Plugin, error) {
	if name == "" {
		return nil, errors.New("no plugin name provided")
	}

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &api.Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"name", bundle.NewCaseInsensitiveRegex(name)}}).Decode(decodedPluginResult)
	if err != nil {
		return nil, err
	}

	return decodedPluginResult, nil

}

func PaginatePlugins(page int) ([]*api.Plugin, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"lastUpdated", -1}})
	if page > 1 {
		findOptions.SetSkip(int64(page*10 - 10))
	}
	findOptions.SetLimit(10)

	collection := session.Client.Database("main").Collection("plugins")

	cur, err := collection.Find(session.Ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	results := []*api.Plugin{}
	defer cur.Close(session.Ctx)
	for cur.Next(session.Ctx) {
		plugin := &api.Plugin{}
		if err = cur.Decode(&plugin); err != nil {
			return nil, err
		}
		results = append(results, plugin)
	}

	return results, nil

}

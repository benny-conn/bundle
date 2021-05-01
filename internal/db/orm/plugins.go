package orm

import (
	"errors"
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PluginsOrm struct{}

func NewPluginsOrm() *PluginsOrm { return &PluginsOrm{} }

func (p *PluginsOrm) Insert(plugin *api.Plugin) error {
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

func (p *PluginsOrm) Update(req *api.UpdatePluginRequest) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")

	updatedPlugin := marshallBsonClean(req.UpdatedPlugin)
	updatedPlugin = append(updatedPlugin, bson.E{"lastUpdated", time.Now().Unix()})

	updateResult, err := collection.UpdateOne(session.Ctx, bson.D{{"name", caseInsensitive(req.Name)}}, bson.D{{"$set", updatedPlugin}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no plugin found")
	}
	return nil

}

func (p *PluginsOrm) Get(req *api.GetPluginRequest) (*api.Plugin, error) {
	if req.Name == "" {
		return nil, errors.New("no plugin name provided")
	}

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("plugins")
	decodedPluginResult := &api.Plugin{}

	err = collection.FindOne(session.Ctx, bson.D{{"name", caseInsensitive(req.Name)}}).Decode(decodedPluginResult)
	if err != nil {
		return nil, err
	}

	return decodedPluginResult, nil

}

func (p *PluginsOrm) Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"lastUpdated", -1}})
	if req.Page > 1 {
		findOptions.SetSkip(int64(req.Page*req.Count - req.Count))
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

	return &api.PaginatePluginsResponse{
		Plugins: results,
	}, nil

}

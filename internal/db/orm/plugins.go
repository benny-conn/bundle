package orm

import (
	"errors"
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrmPlugin struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Author      OrmUser            `bson:"author" json:"author"`
	Version     string             `bson:"version" json:"version"`
	Thumbnail   string             `bson:"thumbnail" json:"thumbnail"`
	LastUpdated primitive.DateTime `bson:"lastUpdated" json:"lastUpdated"`
}

type PluginsOrm struct{}

func NewPluginsOrm() *PluginsOrm { return &PluginsOrm{} }

func (p *PluginsOrm) Insert(plugin *api.Plugin) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("plugins")

	plugin.LastUpdated = time.Now().Unix()

	_, err = collection.InsertOne(session.Ctx, plugin)

	if err != nil {
		return err
	}
	return nil

}

func (p *PluginsOrm) Update(req *api.Plugin) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("plugins")

	updatedPlugin := marshallBsonClean(req)
	updatedPlugin = append(updatedPlugin, bson.E{"lastUpdated", time.Now().Unix()})

	if req.Id == "" {
		plorm := NewPluginsOrm()
		pl, err := plorm.Get(req)
		if err != nil {
			return err
		}
		req.Id = pl.Id
	}

	updateResult, err := collection.UpdateByID(session.Ctx, req.Id, bson.D{{"$set", updatedPlugin}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no plugin found")
	}
	return nil

}

func (p *PluginsOrm) Get(req *api.Plugin) (*api.Plugin, error) {
	if req.Name == "" {
		return nil, errors.New("no plugin name provided")
	}

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("plugins")
	decodedPluginResult := &OrmPlugin{}

	if req.Id == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"name", caseInsensitive(req.Name)}}).Decode(decodedPluginResult)
		if err != nil {
			return nil, err
		}
	} else {
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, err
		}
		err = collection.FindOne(session.Ctx, bson.D{{"_id", id}}).Decode(decodedPluginResult)
		if err != nil {
			return nil, err
		}
	}

	return &api.Plugin{
		Id:          decodedPluginResult.Id.Hex(),
		Name:        decodedPluginResult.Name,
		Description: decodedPluginResult.Description,
		Author: &api.User{
			Id:       decodedPluginResult.Author.Id.Hex(),
			Email:    decodedPluginResult.Author.Email,
			Username: decodedPluginResult.Author.Username,
			Password: decodedPluginResult.Author.Password,
			Tag:      decodedPluginResult.Author.Tag,
			Scopes:   decodedPluginResult.Author.Scopes,
		},
		Version:     decodedPluginResult.Version,
		Thumbnail:   decodedPluginResult.Thumbnail,
		LastUpdated: decodedPluginResult.LastUpdated.Time().Unix(),
	}, nil

}

func (p *PluginsOrm) Paginate(req *api.PaginatePluginsRequest) ([]*api.Plugin, error) {
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
	findOptions.SetLimit(int64(req.Count))

	collection := session.Client.Database("plugins").Collection("plugins")

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

package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Plugin struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Author      User               `bson:"author,omitempty" json:"author"`
	Version     string             `bson:"version,omitempty" json:"version"`
	Thumbnail   string             `bson:"thumbnail,omitempty" json:"thumbnail"`
	LastUpdated primitive.DateTime `bson:"lastUpdated,omitempty" json:"lastUpdated"`
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

	countName, err := collection.CountDocuments(session.Ctx, bson.D{{"name", caseInsensitive(plugin.Name)}})

	if err != nil {
		return err
	}

	if countName > 0 {
		err = errors.New("plugin already exists with given name")
		return err
	}

	insertion := apiToOrmPl(plugin)
	err = validatePluginInsert(insertion)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(session.Ctx, insertion)

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

	update := apiToOrmPl(req)
	err = validatePluginUpdate(update)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateByID(session.Ctx, req.Id, bson.D{{"$set", update}})
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
	decodedPluginResult := &Plugin{}

	get := apiToOrmPl(req)
	err = validatePluginGet(get)
	if err != nil {
		return nil, err
	}

	if get.Id == primitive.NilObjectID {
		err = collection.FindOne(session.Ctx, bson.D{{"name", caseInsensitive(req.Name)}}).Decode(decodedPluginResult)
		if err != nil {
			return nil, err
		}
	} else {
		err = collection.FindOne(session.Ctx, bson.D{{"_id", get.Id}}).Decode(decodedPluginResult)
		if err != nil {
			return nil, err
		}
	}

	return ormToApiPl(*decodedPluginResult), nil

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
		plugin := &Plugin{}
		if err = cur.Decode(&plugin); err != nil {
			return nil, err
		}
		results = append(results, ormToApiPl(*plugin))
	}

	return results, nil

}

func validatePluginUpdate(pl Plugin) error {
	if pl.Id == primitive.NilObjectID {
		return errors.New("id required for update")
	}
	return nil
}

func validatePluginInsert(pl Plugin) error {
	if pl.Name == "" {
		return errors.New("name required for insertion")
	}
	if pl.Version == "" {
		return errors.New("version required for insertion")
	}
	if pl.Author.Id == primitive.NilObjectID {
		return errors.New("author id required for insertion")
	}
	return nil
}

func validatePluginGet(pl Plugin) error {
	if pl.Name == "" && pl.Id == primitive.NilObjectID {
		return errors.New("id or name required for get")
	}
	return nil
}

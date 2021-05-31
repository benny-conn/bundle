package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type changelog struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PluginId primitive.ObjectID `bson:"pluginId,omitempty" json:"pluginId"`
	Version  string             `bson:"version,omitempty" json:"version"`
	Added    []string           `bson:"added,omitempty" json:"added"`
	Removed  []string           `bson:"removed,omitempty" json:"removed"`
	Updated  []string           `bson:"updated,omitempty" json:"updated"`
}
type ChangelogOrm struct{}

func NewChangelogOrm() *ChangelogOrm { return &ChangelogOrm{} }

func (o *ChangelogOrm) Insert(ch *api.Changelog) error {

	mgses, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}
	defer mgses.Cancel()

	collection := mgses.Client.Database("plugins").Collection("changelogs")

	s := apiToOrmChangelog(ch)
	err = validateChangelogInsert(s)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	res, err := collection.InsertOne(mgses.Ctx, s)

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	if res.InsertedID == primitive.NilObjectID {
		err = errors.New("could not insert with new id")
		logger.ErrLog.Print(err.Error())
		return err
	}

	return nil
}

func (o *ChangelogOrm) Get(ch *api.Changelog) (*api.Changelog, error) {
	mgses, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	defer mgses.Cancel()
	collection := mgses.Client.Database("plugins").Collection("changelogs")

	s := apiToOrmChangelog(ch)
	err = validateChangelogGet(s)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	result := collection.FindOne(mgses.Ctx, s)
	if result.Err() != nil {
		logger.ErrLog.Print(result.Err().Error())
		return nil, result.Err()
	}
	final := changelog{}

	err = result.Decode(&final)

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	return ormToApiChangelog(final), nil
}

func (o *ChangelogOrm) GetAll(ch *api.Changelog) (*api.Changelogs, error) {
	mgses, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	defer mgses.Cancel()
	collection := mgses.Client.Database("plugins").Collection("changelogs")

	s := apiToOrmChangelog(ch)
	err = validateChangelogGetAll(s)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	cur, err := collection.Find(mgses.Ctx, s)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	results := []changelog{}

	err = cur.All(mgses.Ctx, &results)

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}

	final := &api.Changelogs{}

	for _, v := range results {
		final.Changelogs = append(final.Changelogs, ormToApiChangelog(v))
	}

	return final, nil
}

func validateChangelogInsert(ch changelog) error {
	if ch.PluginId == primitive.NilObjectID || ch.Version == "" {
		return errors.New("plugin id and version required")
	}
	return nil
}

func validateChangelogGet(ch changelog) error {
	if ch.Id == primitive.NilObjectID {
		if ch.PluginId == primitive.NilObjectID || ch.Version == "" {
			return errors.New("id or plugin id required with version")
		}
	}
	return nil
}

func validateChangelogGetAll(ch changelog) error {
	if ch.PluginId == primitive.NilObjectID {
		return errors.New("plugin id required")
	}

	return nil
}

func apiToOrmChangelog(ch *api.Changelog) changelog {
	if ch == nil {
		return changelog{}
	}
	result := changelog{
		Version: ch.Version,
		Added:   ch.Added,
		Removed: ch.Removed,
		Updated: ch.Updated,
	}

	if ch.Id != "" {
		id, err := primitive.ObjectIDFromHex(ch.Id)
		if err == nil && id != primitive.NilObjectID {
			result.Id = id
		}
	}
	if ch.PluginId != "" {
		id, err := primitive.ObjectIDFromHex(ch.PluginId)
		if err == nil && id != primitive.NilObjectID {
			result.PluginId = id
		}
	}

	return result
}

func ormToApiChangelog(ch changelog) *api.Changelog {
	return &api.Changelog{
		Id:       ch.Id.Hex(),
		PluginId: ch.PluginId.Hex(),
		Version:  ch.Version,
		Added:    ch.Added,
		Removed:  ch.Removed,
		Updated:  ch.Updated,
	}
}

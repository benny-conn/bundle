package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReadmesOrm struct{}

func NewReadmesOrm() *ReadmesOrm { return &ReadmesOrm{} }

func (p *ReadmesOrm) Insert(readme *api.Readme) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("readmes")

	newReadme, err := bson.Marshal(readme)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(session.Ctx, newReadme)

	if err != nil {
		return err
	}
	return nil

}

func (p *ReadmesOrm) Update(req *api.Readme) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("plugins")

	updatedReadme := marshallBsonClean(req)

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return err
	}
	updateResult, err := collection.UpdateByID(session.Ctx, id, bson.D{{"$set", updatedReadme}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no plugin found")
	}
	return nil

}

func (p *ReadmesOrm) Get(req *api.Plugin) (*api.Readme, error) {
	if req.Name == "" {
		return nil, errors.New("no plugin name provided")
	}

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("readmes")
	decodedReadmeResult := &api.Readme{}

	if req.Id == "" {
		plorm := NewPluginsOrm()
		pl, err := plorm.Get(req)
		if err != nil {
			return nil, err
		}
		req.Id = pl.Id
	}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", req.Id}}).Decode(decodedReadmeResult)
	if err != nil {
		return nil, err
	}

	return decodedReadmeResult, nil

}

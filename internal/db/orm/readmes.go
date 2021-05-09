package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Readme struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Plugin primitive.ObjectID `bson:"plugin,omitempty" json:"plugin"`
	Text   string             `bson:"text,omitempty" json:"text"`
}
type ReadmesOrm struct{}

func NewReadmesOrm() *ReadmesOrm { return &ReadmesOrm{} }

func (p *ReadmesOrm) Insert(readme *api.Readme) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("readmes")

	if readme.Plugin == nil {
		return errors.New("plugin not specified")
	}
	var plId primitive.ObjectID
	if readme.Plugin.Id == "" {
		dbpl, err := NewPluginsOrm().Get(readme.Plugin)

		if err != nil {
			return err
		}

		plId, err = primitive.ObjectIDFromHex(dbpl.Id)
		if err != nil {
			return err
		}
	} else {
		plId, err = primitive.ObjectIDFromHex(readme.Plugin.Id)
		if err != nil {
			return err
		}

	}
	count, err := collection.CountDocuments(session.Ctx, bson.D{{"plugin", plId}})

	if err != nil {
		return err
	}

	if count > 0 {
		err = errors.New("plugin already has a readme, please update instead")
		return err
	}

	insert := apiToOrmReadme(readme)

	err = validateReadmeInsert(insert)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(session.Ctx, insert)

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

	collection := session.Client.Database("plugins").Collection("readmes")

	updated := apiToOrmReadme(req)

	err = validateReadmeUpdate(updated)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateByID(session.Ctx, req.Id, bson.D{{"$set", updated}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no plugin found")
	}
	return nil

}

func (p *ReadmesOrm) Get(req *api.Plugin) (*api.Readme, error) {

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("plugins").Collection("readmes")
	decodedReadmeResult := &Readme{}

	var plId primitive.ObjectID
	if req.Id == "" {
		dbpl, err := NewPluginsOrm().Get(req)

		if err != nil {
			return nil, err
		}

		plId, err = primitive.ObjectIDFromHex(dbpl.Id)
		if err != nil {
			return nil, err
		}
	} else {
		plId, err = primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, err
		}

	}

	err = collection.FindOne(session.Ctx, bson.D{{"plugin", plId}}).Decode(decodedReadmeResult)
	if err != nil {
		return nil, err
	}

	return ormToApiReadme(*decodedReadmeResult), nil

}

func validateReadmeGet(pl Plugin) error {
	if pl.Name == "" {
		return errors.New("plugin name required for search")
	}
	return nil
}

func validateReadmeUpdate(rdme Readme) error {
	if rdme.Id == primitive.NilObjectID {
		return errors.New("id required for update")
	}
	return nil
}

func validateReadmeInsert(rdme Readme) error {
	if rdme.Plugin == primitive.NilObjectID {
		return errors.New("valid plugin required for insertion")
	}
	if rdme.Text == "" {
		return errors.New("valid readme required for insertion")
	}
	return nil
}

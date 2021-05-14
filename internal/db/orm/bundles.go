package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bundle struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId  primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	FtpUser string             `bson:"ftpUser,omitempty" json:"ftpUser"`
	FtpPass string             `bson:"ftpPass,omitempty" json:"ftpPass"`
	FtpPort int32              `bson:"ftpPort,omitempty" json:"ftpPort"`
	FtpHost string             `bson:"ftpHost,omitempty" json:"ftpHost"`
	Plugins []string           `bson:"plugins,omitempty" json:"plugins"`
}
type BundlesOrm struct{}

func NewBundlesOrm() *BundlesOrm { return &BundlesOrm{} }

func (o *BundlesOrm) Insert(bu *api.Bundle) error {

	mgses, err := getMongoSession()
	if err != nil {
		return err
	}
	defer mgses.Cancel()

	collection := mgses.Client.Database("users").Collection("bundles")

	s := apiToOrmBundle(bu)
	err = validateBundleInsert(s)
	if err != nil {
		return err
	}

	res, err := collection.InsertOne(mgses.Ctx, s)

	if err != nil {
		return err
	}

	if res.InsertedID == primitive.NilObjectID {
		return errors.New("could not insert with new id")
	}

	return nil
}

func (o *BundlesOrm) Delete(bu *api.Bundle) error {
	mgses, err := getMongoSession()
	if err != nil {
		return err
	}
	defer mgses.Cancel()

	collection := mgses.Client.Database("users").Collection("bundles")

	s := apiToOrmBundle(bu)
	err = validateBundleDelete(s)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(mgses.Ctx, s)

	if err != nil {
		return err
	}
	return nil
}

func (o *BundlesOrm) Update(bu *api.Bundle) error {
	mgses, err := getMongoSession()
	if err != nil {
		return err
	}
	defer mgses.Cancel()

	collection := mgses.Client.Database("users").Collection("bundles")

	s := apiToOrmBundle(bu)
	err = validateBundleUpdate(s)
	if err != nil {
		return err
	}

	update, err := collection.UpdateByID(mgses.Ctx, bu.Id, bson.D{{"$set", bu}})
	if update.MatchedCount < 1 {
		return errors.New("could not find a session to updae last retrieved time")
	}

	if err != nil {
		return err
	}

	return nil
}

func (o *BundlesOrm) Get(bu *api.Bundle) (*api.Bundle, error) {
	mgses, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer mgses.Cancel()
	collection := mgses.Client.Database("users").Collection("bundles")

	s := apiToOrmBundle(bu)
	err = validateBundleGet(s)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(mgses.Ctx, s)
	if result.Err() != nil {
		return nil, result.Err()
	}
	final := bundle{}

	err = result.Decode(&final)

	if err != nil {
		return nil, err
	}

	return ormToApiBundle(final), nil
}

func validateBundleInsert(bu bundle) error {
	if bu.UserId == primitive.NilObjectID {
		return errors.New("user id required")
	}
	return nil
}

func validateBundleGet(bu bundle) error {
	if bu.Id == primitive.NilObjectID && bu.UserId == primitive.NilObjectID {
		return errors.New("id or user id required")
	}
	return nil
}

func validateBundleUpdate(bu bundle) error {
	if bu.Id == primitive.NilObjectID {
		return errors.New("id required")
	}
	return nil
}

func validateBundleDelete(bu bundle) error {
	if bu.Id == primitive.NilObjectID {
		return errors.New("id required")
	}
	return nil
}

func apiToOrmBundle(bu *api.Bundle) bundle {
	if bu == nil {
		return bundle{}
	}
	result := bundle{
		FtpUser: bu.FtpUser,
		FtpPort: bu.FtpPort,
		FtpHost: bu.FtpHost,
		Plugins: bu.Plugins,
	}

	if bu.Id != "" {
		id, err := primitive.ObjectIDFromHex(bu.Id)
		if err == nil && id != primitive.NilObjectID {
			result.Id = id
		}
	}
	if bu.UserId != "" {
		id, err := primitive.ObjectIDFromHex(bu.UserId)
		if err == nil && id != primitive.NilObjectID {
			result.UserId = id
		}
	}

	return result
}

func ormToApiBundle(bu bundle) *api.Bundle {
	return &api.Bundle{
		Id:      bu.Id.Hex(),
		UserId:  bu.UserId.Hex(),
		FtpUser: bu.FtpUser,
		FtpPort: bu.FtpPort,
		FtpHost: bu.FtpHost,
		Plugins: bu.Plugins,
	}
}

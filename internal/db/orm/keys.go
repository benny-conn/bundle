package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type key struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"userId,omitempty" json:"plugin"`
}

type KeysOrm struct{}

func NewKeysOrm() *KeysOrm { return &KeysOrm{} }

func (p *KeysOrm) Insert(k *api.Key) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("keys")

	insert := apiToOrmKey(k)

	err = validateKeyInsert(insert)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(session.Ctx, insert)

	if err != nil {
		return err
	}
	return nil

}

func (p *KeysOrm) Get(req *api.Key) (*api.Key, error) {

	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("keys")
	decodedReadmeResult := &key{}

	get := apiToOrmKey(req)
	err = validateKeyGet(get)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(session.Ctx, bson.D{{"userId", req.UserId}}).Decode(decodedReadmeResult)
	if err != nil {
		return nil, err
	}

	return ormToApiKeys(*decodedReadmeResult), nil

}

func (p *KeysOrm) Delete(k *api.Key) error {

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("keys")

	delete := apiToOrmKey(k)

	err = validateKeyDelete(delete)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(session.Ctx, delete)

	if err != nil {
		return err
	}
	return nil

}

func validateKeyGet(k key) error {
	if k.UserId == primitive.NilObjectID {
		return errors.New("userid required for get")
	}
	return nil
}

func validateKeyInsert(k key) error {
	if k.UserId == primitive.NilObjectID {
		return errors.New("valid userId required for insertion")
	}
	return nil
}
func validateKeyDelete(k key) error {
	if k.Id == primitive.NilObjectID && k.UserId == primitive.NilObjectID {
		return errors.New("valid id or userId required for deletion")
	}

	return nil
}

func ormToApiKeys(k key) *api.Key {
	r := &api.Key{
		Id:     k.Id.Hex(),
		UserId: k.UserId.Hex(),
	}

	return r
}

func apiToOrmKey(k *api.Key) key {
	if k == nil {
		return key{}
	}
	result := key{}
	id, err := primitive.ObjectIDFromHex(k.Id)
	if id != primitive.NilObjectID && err == nil {
		result.Id = id
	}
	userId, err := primitive.ObjectIDFromHex(k.UserId)
	if userId != primitive.NilObjectID && err == nil {
		result.UserId = userId
	}

	return result
}

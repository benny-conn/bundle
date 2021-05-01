package orm

import (
	"context"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func getMongoSession() (*Mongo, error) {
	mg := &Mongo{}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MongoURL")))
	mg.Cancel = cancel
	mg.Client = client
	mg.Ctx = ctx
	if err != nil {
		return mg, err
	}
	return mg, nil
}

func marshallBsonClean(val interface{}) bson.D {

	bs, _ := bson.Marshal(val)

	ogBson := bson.D{}
	bson.Unmarshal(bs, &ogBson)

	newVal := removeZeroOrNilValues(ogBson)
	return newVal
}

func removeZeroOrNilValues(val bson.D) bson.D {
	b := val
	for i, v := range b {
		refVal := reflect.ValueOf(v.Value)
		if !refVal.IsValid() {
			b = remove(b, i)
			return removeZeroOrNilValues(b)
		}
		if refVal.IsZero() {
			b = remove(b, i)
			return removeZeroOrNilValues(b)
		}

	}
	return b
}

func remove(s bson.D, i int) bson.D {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func caseInsensitive(value string) primitive.Regex {
	return primitive.Regex{Pattern: value, Options: "i"}
}
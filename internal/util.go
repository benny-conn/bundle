package internal

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func NewCaseInsensitiveRegex(value string) primitive.Regex {
	return primitive.Regex{Pattern: value, Options: "i"}
}

func GetMongoSession() (*Mongo, error) {
	mg := &Mongo{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoURL))
	mg.Cancel = cancel
	mg.Client = client
	mg.Ctx = ctx
	if err != nil {
		return mg, err
	}
	return mg, nil
}

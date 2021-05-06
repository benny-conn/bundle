package orm

import (
	"errors"
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId        primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	LastRetrieved primitive.DateTime `bson:"lastRetrieved,omitempty" json:"lastRetrieved"`
	CreatedAt     primitive.DateTime `bson:"createdAt,omitempty" json:"createdAt"`
}
type SessionsOrm struct{}

func NewSessionsOrm() *SessionsOrm { return &SessionsOrm{} }

func (o *SessionsOrm) Insert(ses *api.Session) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("sessions")

	s := apiToOrmSession(ses)
	err = validateSesInsert(s)
	if err != nil {
		return err
	}
	s.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	s.LastRetrieved = primitive.NewDateTimeFromTime(time.Now())

	res, err := collection.InsertOne(session.Ctx, s)

	if err != nil {
		return err
	}

	if res.InsertedID == primitive.NilObjectID {
		return errors.New("could not insert with new id")
	}

	return nil
}

func (o *SessionsOrm) Delete(ses *api.Session) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("sessions")

	s := apiToOrmSession(ses)
	err = validateSesDelete(s)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(session.Ctx, s)

	if err != nil {
		return err
	}
	return nil
}

func (o *SessionsOrm) Get(ses *api.Session) (*api.Session, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()
	collection := session.Client.Database("users").Collection("sessions")

	s := apiToOrmSession(ses)
	err = validateSesGet(s)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(session.Ctx, s)
	if result.Err() != nil {
		return nil, result.Err()
	}
	final := Session{}

	err = result.Decode(&final)

	if err != nil {
		return nil, err
	}

	upSes := final
	upSes.LastRetrieved = primitive.NewDateTimeFromTime(time.Now())

	update, err := collection.UpdateByID(session.Ctx, upSes.Id, bson.D{{"$set", upSes}})
	if update.MatchedCount < 1 {
		return nil, errors.New("could not find a session to updae last retrieved time")
	}

	if err != nil {
		return nil, err
	}

	return ormToApiSession(final), nil
}

func validateSesInsert(ses Session) error {
	if ses.UserId == primitive.NilObjectID {
		return errors.New("user id required")
	}
	return nil
}

func validateSesGet(ses Session) error {
	if ses.Id == primitive.NilObjectID && ses.UserId == primitive.NilObjectID {
		return errors.New("id or user id required")
	}
	return nil
}

func validateSesDelete(ses Session) error {
	if ses.Id == primitive.NilObjectID {
		return errors.New("id required")
	}
	return nil
}

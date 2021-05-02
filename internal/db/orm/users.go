package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OrmUser struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Tag      string             `bson:"tag" json:"tag"`
	Scopes   []string           `bson:"scopes" json:"scopes"`
}

type UsersOrm struct{}

func NewUsersOrm() *UsersOrm { return &UsersOrm{} }

func (u *UsersOrm) Insert(user *api.User) error {
	isValid := isUserValid(user)

	if !isValid {
		return errors.New("invalid user")
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(bcryptPass)

	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", caseInsensitive(user.Username)}})

	if err != nil {
		return err
	}

	if countUserName > 0 {
		err = errors.New("user already exists with given username")
		return err
	}

	countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", caseInsensitive(user.Email)}})

	if err != nil {
		return err
	}

	if countEmail > 0 {
		err = errors.New("user already exists with given email")
		return err
	}

	_, err = collection.InsertOne(session.Ctx, user)

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersOrm) Get(req *api.User) (*api.User, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	decodedUser := &OrmUser{}

	switch {
	case req.Id != "":
		id, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, err
		}
		err = collection.FindOne(session.Ctx, bson.D{{"_id", id}}).Decode(decodedUser)
	case req.Email == "":
		err = collection.FindOne(session.Ctx, bson.D{{"username", req.Username}}).Decode(decodedUser)
	case req.Username == "":
		err = collection.FindOne(session.Ctx, bson.D{{"email", caseInsensitive(req.Email)}}).Decode(decodedUser)
	default:
		err = collection.FindOne(session.Ctx, bson.D{{"username", req.Username}, {"email", caseInsensitive(req.Email)}}).Decode(decodedUser)
	}

	if err != nil {
		return nil, err
	}

	return &api.User{
		Id:       decodedUser.Id.Hex(),
		Email:    decodedUser.Email,
		Username: decodedUser.Username,
		Password: decodedUser.Password,
		Tag:      decodedUser.Tag,
		Scopes:   decodedUser.Scopes,
	}, nil
}

func (u *UsersOrm) Update(req *api.User) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	updatedUser := marshallBsonClean(req)

	if req.Id == "" {
		user, err := u.Get(req)
		if err != nil {
			return err
		}
		req.Id = user.Id
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateByID(session.Ctx, id, bson.D{{"$set", updatedUser}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no user found")
	}
	return nil
}

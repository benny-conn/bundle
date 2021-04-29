package orm

import (
	"errors"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UsersOrm struct{}

func NewUsersOrm() internal.UserService { return &UsersOrm{} }

func (u *UsersOrm) Insert(user *api.User) error {
	isValid := internal.IsUserValid(user)

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

	collection := session.Client.Database("main").Collection("users")

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

	userToInsert := marshallBsonClean(user)

	_, err = collection.InsertOne(session.Ctx, userToInsert)

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersOrm) Get(req *api.GetUserRequest) (*api.User, error) {
	session, err := getMongoSession()
	if err != nil {
		return nil, err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	decodedUser := &api.User{}

	if req.Email == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"username", req.Username}}).Decode(decodedUser)
	} else if req.Username == "" {
		err = collection.FindOne(session.Ctx, bson.D{{"email", caseInsensitive(req.Email)}}).Decode(decodedUser)
	} else {
		err = collection.FindOne(session.Ctx, bson.D{{"username", req.Username}, {"email", caseInsensitive(req.Email)}}).Decode(decodedUser)
	}
	if err != nil {
		return nil, err
	}

	return decodedUser, nil
}

func (u *UsersOrm) Update(req *api.UpdateUserRequest) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("main").Collection("users")

	updatedUser := marshallBsonClean(req.UpdatedUser)

	updateResult, err := collection.UpdateOne(session.Ctx, bson.D{{"username", caseInsensitive(req.Username)}}, bson.D{{"$set", updatedUser}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no user found")
	}
	return nil
}

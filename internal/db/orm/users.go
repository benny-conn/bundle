package orm

import (
	"errors"
	"regexp"

	"github.com/bennycio/bundle/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username,omitempty" json:"username"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Password string             `bson:"password,omitempty" json:"password"`
	Tag      string             `bson:"tag,omitempty" json:"tag"`
	Scopes   []string           `bson:"scopes,omitempty" json:"scopes"`
}

type UsersOrm struct{}

func NewUsersOrm() *UsersOrm { return &UsersOrm{} }

func (u *UsersOrm) Insert(user *api.User) error {

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

	insertion := apiToOrmUser(user)
	err = validateUserInsert(insertion)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(session.Ctx, insertion)

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

	decodedUser := &User{}
	get := apiToOrmUser(req)
	err = validateUserGet(get)
	if err != nil {
		return nil, err
	}
	switch {
	case get.Id != primitive.NilObjectID:
		res := collection.FindOne(session.Ctx, bson.D{{"_id", get.Id}})
		if res.Err() != nil {
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	case get.Email == "":
		res := collection.FindOne(session.Ctx, bson.D{{"username", get.Username}})
		if res.Err() != nil {
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	case get.Username == "":
		res := collection.FindOne(session.Ctx, bson.D{{"email", caseInsensitive(get.Email)}})
		if res.Err() != nil {
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	default:
		res := collection.FindOne(session.Ctx, bson.D{{"username", get.Username}, {"email", caseInsensitive(get.Email)}})
		if res.Err() != nil {
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	}

	return ormToApiUser(*decodedUser), nil
}

func (u *UsersOrm) Update(req *api.User) error {
	session, err := getMongoSession()
	if err != nil {
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	update := apiToOrmUser(req)
	err = validateUserUpdate(update)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateByID(session.Ctx, update.Id, bson.D{{"$set", update}})
	if err != nil {
		return err
	}
	if updateResult.MatchedCount < 1 {
		return errors.New("no user found")
	}
	return nil
}

func validateUserGet(user User) error {

	if user.Id == primitive.NilObjectID && user.Email == "" && user.Username == "" {
		return errors.New("id, email, or username is required for get")
	}
	return nil
}

func validateUserInsert(user User) error {

	if user.Username == "" {
		return errors.New("username required for insert")
	}

	if user.Password == "" {
		return errors.New("password required for insert")
	}

	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.Email != "" {
		if len(user.Email) > 254 || !rxEmail.MatchString(user.Email) {
			return errors.New("invalid email")
		}
	}

	return nil

}

func validateUserUpdate(user User) error {
	if user.Id == primitive.NilObjectID {
		return errors.New("id required for update")
	}
	return nil
}

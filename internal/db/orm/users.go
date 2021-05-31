package orm

import (
	"errors"
	"regexp"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username,omitempty" json:"username"`
	Email     string             `bson:"email,omitempty" json:"email"`
	Password  string             `bson:"password,omitempty" json:"password"`
	Tag       string             `bson:"tag,omitempty" json:"tag"`
	Scopes    []string           `bson:"scopes,omitempty" json:"scopes"`
	Thumbnail string             `bson:"thumbnail,omitempty" json:"thumbnail"`
	StripeId  string             `bson:"stripeId,omitempty" json:"stripeId"`
}

type UsersOrm struct{}

func NewUsersOrm() *UsersOrm { return &UsersOrm{} }

func (u *UsersOrm) Insert(us *api.User) error {

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	us.Password = string(bcryptPass)

	session, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	countUserName, err := collection.CountDocuments(session.Ctx, bson.D{{"username", caseInsensitive(us.Username)}})

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	if countUserName > 0 {
		err = errors.New("user already exists with given username")
		logger.ErrLog.Print(err.Error())
		return err
	}

	countEmail, err := collection.CountDocuments(session.Ctx, bson.D{{"email", caseInsensitive(us.Email)}})

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	if countEmail > 0 {
		err = errors.New("user already exists with given email")
		logger.ErrLog.Print(err.Error())
		return err
	}

	insertion := apiToOrmUser(us)
	err = validateUserInsert(insertion)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	_, err = collection.InsertOne(session.Ctx, insertion)

	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}
	return nil
}

func (u *UsersOrm) Get(req *api.User) (*api.User, error) {

	session, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	defer session.Cancel()
	collection := session.Client.Database("users").Collection("users")

	decodedUser := &user{}
	get := apiToOrmUser(req)
	err = validateUserGet(get)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return nil, err
	}
	switch {
	case get.Id != primitive.NilObjectID:
		res := collection.FindOne(session.Ctx, bson.D{{"_id", get.Id}})
		if res.Err() != nil {
			logger.ErrLog.Print(res.Err().Error())
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	case get.Email == "":
		res := collection.FindOne(session.Ctx, bson.D{{"username", get.Username}})
		if res.Err() != nil {
			logger.ErrLog.Print(res.Err().Error())
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	case get.Username == "":
		res := collection.FindOne(session.Ctx, bson.D{{"email", caseInsensitive(get.Email)}})
		if res.Err() != nil {
			logger.ErrLog.Print(res.Err().Error())
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	default:
		res := collection.FindOne(session.Ctx, bson.D{{"username", get.Username}, {"email", caseInsensitive(get.Email)}})
		if res.Err() != nil {
			logger.ErrLog.Print(res.Err().Error())
			return nil, res.Err()
		}
		res.Decode(decodedUser)
	}

	return ormToApiUser(*decodedUser), nil
}

func (u *UsersOrm) Update(req *api.User) error {
	session, err := getMongoSession()
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}
	defer session.Cancel()

	collection := session.Client.Database("users").Collection("users")

	update := apiToOrmUser(req)
	err = validateUserUpdate(update)
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}

	updateResult, err := collection.UpdateByID(session.Ctx, update.Id, bson.D{{"$set", update}})
	if err != nil {
		logger.ErrLog.Print(err.Error())
		return err
	}
	if updateResult.MatchedCount < 1 {
		err = errors.New("no user found")
		logger.ErrLog.Print(err.Error())
		return err
	}
	return nil
}

func validateUserGet(us user) error {

	if us.Id == primitive.NilObjectID && us.Email == "" && us.Username == "" {
		return errors.New("id, email, or username is required for get")
	}
	return nil
}

func validateUserInsert(us user) error {

	if us.Username == "" {
		return errors.New("username required for insert")
	}

	if us.Password == "" {
		return errors.New("password required for insert")
	}

	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if us.Email != "" {
		if len(us.Email) > 254 || !rxEmail.MatchString(us.Email) {
			return errors.New("invalid email")
		}
	}

	return nil

}

func validateUserUpdate(us user) error {
	if us.Id == primitive.NilObjectID {
		return errors.New("id required for update")
	}
	return nil
}

func apiToOrmUser(us *api.User) user {
	if us == nil {
		return user{}
	}
	result := user{
		Username:  us.Username,
		Email:     us.Email,
		Password:  us.Password,
		Tag:       us.Tag,
		Scopes:    us.Scopes,
		Thumbnail: us.Thumbnail,
		StripeId:  us.StripeId,
	}
	userID, err := primitive.ObjectIDFromHex(us.Id)
	if userID != primitive.NilObjectID && err == nil {
		result.Id = userID
	}
	return result
}

func ormToApiUser(us user) *api.User {
	return &api.User{
		Id:        us.Id.Hex(),
		Username:  us.Username,
		Email:     us.Email,
		Password:  us.Password,
		Tag:       us.Tag,
		Scopes:    us.Scopes,
		Thumbnail: us.Thumbnail,
		StripeId:  us.StripeId,
	}
}

package orm

import (
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ormToApiPl(pl Plugin) *api.Plugin {
	return &api.Plugin{
		Id:          pl.Id.Hex(),
		Name:        pl.Name,
		Description: pl.Description,
		Author:      ormToApiUser(pl.Author),
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: pl.LastUpdated.Time().Unix(),
	}
}

func apiToOrmPl(pl *api.Plugin) Plugin {

	pluginID, err := primitive.ObjectIDFromHex(pl.Id)
	if err != nil {
		pluginID = primitive.NilObjectID
	}
	authorID, err := primitive.ObjectIDFromHex(pl.Author.Id)
	if err != nil {
		authorID = primitive.NilObjectID
	}
	var lastUpdated primitive.DateTime
	if pl.LastUpdated != 0 {
		lastUpdated = primitive.DateTime(pl.LastUpdated)
	} else {
		lastUpdated = primitive.DateTime(time.Now().Unix())
	}
	return Plugin{
		Id:          pluginID,
		Name:        pl.Name,
		Description: pl.Description,
		Author: User{
			Id:       authorID,
			Email:    pl.Author.Email,
			Username: pl.Author.Username,
			Password: pl.Author.Password,
			Tag:      pl.Author.Tag,
			Scopes:   pl.Author.Scopes,
		},
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: lastUpdated,
	}

}

func apiToOrmUser(user *api.User) User {
	userID, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		userID = primitive.NilObjectID
	}
	return User{
		Id:       userID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Tag:      user.Tag,
		Scopes:   user.Scopes,
	}
}

func ormToApiUser(user User) *api.User {
	return &api.User{
		Id:       user.Id.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Tag:      user.Tag,
		Scopes:   user.Scopes,
	}
}

func ormToApiReadme(rdme Readme) *api.Readme {
	return &api.Readme{
		Id:     rdme.Id.Hex(),
		Plugin: ormToApiPl(rdme.Plugin),
		Text:   rdme.Text,
	}
}

func apiToOrmReadme(rdme *api.Readme) Readme {
	id, err := primitive.ObjectIDFromHex(rdme.Id)
	if err != nil {
		id = primitive.NilObjectID
	}
	return Readme{
		Id:     id,
		Plugin: apiToOrmPl(rdme.Plugin),
		Text:   rdme.Text,
	}
}

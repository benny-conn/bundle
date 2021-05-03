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

	if pl == nil {
		return Plugin{}
	}

	var lastUpdated primitive.DateTime
	if pl.LastUpdated != 0 {
		lastUpdated = primitive.DateTime(pl.LastUpdated)
	} else {
		lastUpdated = primitive.DateTime(time.Now().Unix())
	}
	result := Plugin{
		Name:        pl.Name,
		Description: pl.Description,
		Author:      apiToOrmUser(pl.Author),
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: lastUpdated,
	}
	pluginID, err := primitive.ObjectIDFromHex(pl.Id)
	if pluginID != primitive.NilObjectID && err == nil {
		result.Id = pluginID
	}

	return result

}

func apiToOrmUser(user *api.User) User {
	if user == nil {
		return User{}
	}
	result := User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Tag:      user.Tag,
		Scopes:   user.Scopes,
	}
	userID, err := primitive.ObjectIDFromHex(user.Id)
	if userID != primitive.NilObjectID && err == nil {
		result.Id = userID
	}
	return result
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
	if rdme == nil {
		return Readme{}
	}
	result := Readme{
		Plugin: apiToOrmPl(rdme.Plugin),
		Text:   rdme.Text,
	}
	id, err := primitive.ObjectIDFromHex(rdme.Id)
	if id != primitive.NilObjectID && err == nil {
		result.Id = id
	}
	return result
}

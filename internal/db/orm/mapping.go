package orm

import (
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ormToApiPl(pl Plugin) *api.Plugin {
	p := &api.Plugin{
		Id:          pl.Id.Hex(),
		Name:        pl.Name,
		Description: pl.Description,
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: pl.LastUpdated.Time().Unix(),
	}
	a, err := NewUsersOrm().Get(&api.User{Id: pl.Author.Hex()})
	if err == nil {
		p.Author = a
	}
	return p
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
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: lastUpdated,
	}
	pluginID, err := primitive.ObjectIDFromHex(pl.Id)
	if pluginID != primitive.NilObjectID && err == nil {
		result.Id = pluginID
	}
	if pl.Author != nil {
		authorID, err := primitive.ObjectIDFromHex(pl.Author.Id)
		if authorID != primitive.NilObjectID && err == nil {
			result.Author = authorID
		}
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
	r := &api.Readme{
		Id:   rdme.Id.Hex(),
		Text: rdme.Text,
	}
	pl, err := NewPluginsOrm().Get(&api.Plugin{Id: rdme.Plugin.Hex()})
	if err == nil {
		r.Plugin = pl
	}
	return r
}

func apiToOrmReadme(rdme *api.Readme) Readme {
	if rdme == nil {
		return Readme{}
	}
	result := Readme{
		Text: rdme.Text,
	}
	id, err := primitive.ObjectIDFromHex(rdme.Id)
	if id != primitive.NilObjectID && err == nil {
		result.Id = id
	}
	if rdme.Plugin != nil {
		pl, err := primitive.ObjectIDFromHex(rdme.Plugin.Id)
		if pl != primitive.NilObjectID && err == nil {
			result.Plugin = pl
		}
	}
	return result
}

func apiToOrmSession(ses *api.Session) Session {
	if ses == nil {
		return Session{}
	}
	result := Session{}
	if ses.Id != "" {
		id, err := primitive.ObjectIDFromHex(ses.Id)
		if err == nil && id != primitive.NilObjectID {
			result.Id = id
		}
	}
	if ses.UserId != "" {
		id, err := primitive.ObjectIDFromHex(ses.UserId)
		if err == nil && id != primitive.NilObjectID {
			result.UserId = id
		}
	}
	if ses.LastRetrieved != 0 {
		result.LastRetrieved = primitive.DateTime(ses.LastRetrieved)
	}
	if ses.CreatedAt != 0 {
		result.CreatedAt = primitive.DateTime(ses.CreatedAt)
	}
	return result
}

func ormToApiSession(ses Session) *api.Session {
	return &api.Session{
		Id:            ses.Id.Hex(),
		UserId:        ses.UserId.Hex(),
		LastRetrieved: ses.LastRetrieved.Time().Unix(),
		CreatedAt:     ses.CreatedAt.Time().Unix(),
	}
}

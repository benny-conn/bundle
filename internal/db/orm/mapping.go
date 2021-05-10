package orm

import (
	"time"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ormToApiPl(pl plugin) *api.Plugin {
	p := &api.Plugin{
		Id:          pl.Id.Hex(),
		Name:        pl.Name,
		Description: pl.Description,
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		Category:    api.Category(pl.Category),
		Downloads:   pl.Downloads,
		IsPremium:   pl.IsPremium,
		Premium: &api.Premium{
			Price:     pl.Premium.Price,
			Purchases: pl.Premium.Purchases,
		},
		LastUpdated: pl.LastUpdated.Time().Unix(),
	}
	a, err := NewUsersOrm().Get(&api.User{Id: pl.Author.Hex()})
	if err == nil {
		p.Author = a
	}
	return p
}

func apiToOrmPl(pl *api.Plugin) plugin {

	if pl == nil {
		return plugin{}
	}

	var lastUpdated primitive.DateTime
	if pl.LastUpdated != 0 {
		lastUpdated = primitive.DateTime(pl.LastUpdated)
	} else {
		lastUpdated = primitive.DateTime(time.Now().Unix())
	}
	result := plugin{
		Name:        pl.Name,
		Description: pl.Description,
		Version:     pl.Version,
		Thumbnail:   pl.Thumbnail,
		LastUpdated: lastUpdated,
		Category:    category(pl.Category),
		Downloads:   pl.Downloads,
		IsPremium:   pl.IsPremium,
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
	if pl.Premium != nil {
		result.Premium = premium{
			Price:     result.Premium.Price,
			Purchases: result.Premium.Purchases,
		}
	}

	return result

}

func apiToOrmUser(us *api.User) user {
	if us == nil {
		return user{}
	}
	result := user{
		Username: us.Username,
		Email:    us.Email,
		Password: us.Password,
		Tag:      us.Tag,
		Scopes:   us.Scopes,
	}
	userID, err := primitive.ObjectIDFromHex(us.Id)
	if userID != primitive.NilObjectID && err == nil {
		result.Id = userID
	}
	return result
}

func ormToApiUser(us user) *api.User {
	return &api.User{
		Id:       us.Id.Hex(),
		Username: us.Username,
		Email:    us.Email,
		Password: us.Password,
		Tag:      us.Tag,
		Scopes:   us.Scopes,
	}
}

func ormToApiReadme(rdme readme) *api.Readme {
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

func apiToOrmReadme(rdme *api.Readme) readme {
	if rdme == nil {
		return readme{}
	}
	result := readme{
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

func apiToOrmSession(ses *api.Session) session {
	if ses == nil {
		return session{}
	}
	result := session{}
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

func ormToApiSession(ses session) *api.Session {
	return &api.Session{
		Id:            ses.Id.Hex(),
		UserId:        ses.UserId.Hex(),
		LastRetrieved: ses.LastRetrieved.Time().Unix(),
		CreatedAt:     ses.CreatedAt.Time().Unix(),
	}
}

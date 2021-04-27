package web

import (
	"errors"
	"net/http"
	"strings"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/internal/storage"
)

func getProfileFromCookie(r *http.Request) (bundle.User, error) {

	c, err := r.Cookie("access_token")
	if err != nil {
		return bundle.User{}, err
	}
	user, err := auth.GetUserFromToken(c.Value)
	if err != nil {
		return bundle.User{}, err
	}
	return user, nil

}

func validateAndReturnPlugin(plugin bundle.Plugin) (bundle.Plugin, error) {
	isReadme := plugin.Version == "README"
	dbPlugin, err := storage.GetPlugin(plugin.Name)

	if err == nil {
		isUserPluginAuthor := strings.EqualFold(dbPlugin.Author, plugin.Author)

		if isUserPluginAuthor {
			if isReadme {
				dbPlugin.Version = plugin.Version
			} else {
				err = storage.UpdatePlugin(plugin.Name, plugin)
				if err != nil {
					return bundle.Plugin{}, err
				}
			}
		} else {
			return bundle.Plugin{}, err
		}
	} else {
		if isReadme {
			err = errors.New("no plugin to attach readme to")
			return bundle.Plugin{}, err
		}
		err = storage.InsertPlugin(plugin)
		if err != nil {
			return bundle.Plugin{}, err
		}
		return plugin, nil
	}
	return dbPlugin, nil
}

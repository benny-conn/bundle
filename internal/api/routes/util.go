package routes

import (
	"errors"
	"strings"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
)

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

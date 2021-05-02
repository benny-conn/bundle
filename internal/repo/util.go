package repo

import (
	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

// Make sure plugin has an ID and an author ID.
// If not, grab those fields from db and write them
// to the plugin. Will return error if plugin is not
// in db to begin with.
func cleanPlugin(plugin *api.Plugin) error {

	gs := gate.NewGateService("", "")

	dbpl, err := gs.GetPlugin(plugin)
	if err != nil {
		return err
	}

	if plugin.Id == "" {
		plugin.Id = dbpl.Id
	}
	if plugin.Author == "" {
		plugin.Author = dbpl.Author
	}
	return nil

}

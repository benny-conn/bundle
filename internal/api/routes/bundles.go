package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	pb "github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
)

func BundleHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")

		plugin, err := getPlugin(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if version != "latest" {
			plugin.Version = version
		}

		bs, err := storage.DownloadFromRepo(plugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		bundle.WriteResponse(w, string(bs), http.StatusOK)
	}

	if r.Method == http.MethodPost {

		plugin := r.Header.Get("Plugin")

		reqPlugin := &pb.Plugin{}

		err := json.Unmarshal([]byte(plugin), reqPlugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		finalPlugin, err := updateOrInsertPlugin(reqPlugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uploadLocation, err := storage.UploadToRepo(finalPlugin, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to", uploadLocation)
	}

}

func updateOrInsertPlugin(plugin *pb.Plugin) (*pb.Plugin, error) {
	isReadme := plugin.Version == "README"
	dbPlugin, err := getPlugin(plugin.Name)

	if err == nil {
		isUserPluginAuthor := strings.EqualFold(dbPlugin.Author, plugin.Author)

		if isUserPluginAuthor {
			if isReadme {
				dbPlugin.Version = plugin.Version
			} else {
				err = updatePlugin(plugin.Name, plugin)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		if isReadme {
			err = errors.New("no plugin to attach readme to")
			return nil, err
		}
		err = insertPlugin(plugin)
		if err != nil {
			return nil, err
		}
		return plugin, nil
	}
	return dbPlugin, nil
}

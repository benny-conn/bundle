package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

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

		plugin, err := storage.GetPlugin(name)
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

		reqPlugin := &bundle.Plugin{}

		err := json.Unmarshal([]byte(plugin), reqPlugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		finalPlugin, err := validateAndReturnPlugin(*reqPlugin)
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

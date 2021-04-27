package routes

import (
	"encoding/json"
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
)

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("plugin")

		plugin, err := storage.GetPlugin(pluginName)
		if err != nil {
			panic(err)
		}

		asJSON, err := json.Marshal(plugin)
		if err != nil {
			panic(err)
		}

		bundle.WriteResponse(w, string(asJSON), http.StatusOK)
	}
}

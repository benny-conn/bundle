package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
)

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("name")
		page := r.FormValue("page")

		if pluginName != "" {
			plugin, err := wrapper.GetPlugin(pluginName)
			if err != nil {
				panic(err)
			}

			asJSON, err := json.Marshal(plugin)
			if err != nil {
				panic(err)
			}

			bundle.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		} else if page != "" {
			convPage, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			plugins, err := wrapper.PaginatePlugins(convPage)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			asJSON, err := json.Marshal(plugins)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bundle.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		}
	}
}

package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
)

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

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

			internal.WriteResponse(w, string(asJSON), http.StatusOK)
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
			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		}
	}
}

func UsersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		bs, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser := &api.User{}
		err = json.Unmarshal(bs, newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = wrapper.InsertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

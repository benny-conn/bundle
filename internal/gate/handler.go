package gate

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
)

func usersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	client := newUserClient("", "")

	switch req.Method {
	case http.MethodGet:
		req.ParseForm()

		userName := req.FormValue("username")
		email := req.FormValue("email")

		r := &api.User{
			Username: userName,
			Email:    email,
		}
		user, err := client.Get(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bs, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bs)
	case http.MethodPost:
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

		err = client.Insert(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
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

		err = client.Update(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func pluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	client := newPluginClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("name")
		page := r.FormValue("page")
		count := r.FormValue("count")

		if pluginName != "" {

			req := &api.Plugin{
				Name: pluginName,
			}
			plugin, err := client.Get(req)
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
			convCount, err := strconv.Atoi(count)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			req := &api.PaginatePluginsRequest{
				Page:  int32(convPage),
				Count: int32(convCount),
			}
			plugins, err := client.Paginate(req)
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
	case http.MethodPost:
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pl := &api.Plugin{}

		err = json.Unmarshal(bs, pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Insert(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Plugin{}

		err = json.Unmarshal(bs, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Update(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func readmesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func repoPluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// TODO
}

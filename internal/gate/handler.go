package gate

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/repo"
)

func usersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	client := newUserClient("", "")

	switch req.Method {
	case http.MethodGet:
		req.ParseForm()

		userName := req.FormValue("username")
		email := req.FormValue("email")
		id := req.FormValue("id")

		r := &api.User{
			Id:       id,
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
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			asJSON, err := json.Marshal(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
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
	client := newReadmeClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pluginName := r.FormValue("name")
		id := r.FormValue("id")

		req := &api.Plugin{
			Id:   id,
			Name: pluginName,
		}
		rdme, err := client.Get(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		asJSON, err := json.Marshal(rdme)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		internal.WriteResponse(w, string(asJSON), http.StatusOK)
		return

	case http.MethodPost:
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pl := &api.Readme{}

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

		req := &api.Readme{}

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

func repoPluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	repo := repo.NewRepoService("", "")
	gs := NewGateService("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pluginName := r.FormValue("name")
		id := r.FormValue("id")

		req := &api.Plugin{
			Id:   id,
			Name: pluginName,
		}
		pl, err := repo.DownloadPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(pl)
		return

	case http.MethodPost:

		userJSON := r.Header.Get("User")
		pluginJSON := r.Header.Get("Resource")

		if userJSON == "" || pluginJSON == "" {
			http.Error(w, "incomplete headers", http.StatusBadRequest)
			return
		}
		user := &api.User{}
		plugin := &api.Plugin{}

		err := json.Unmarshal([]byte(userJSON), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal([]byte(pluginJSON), plugin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := gs.GetUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		plugin.Author = dbUser
		dbPlugin, err := gs.GetPlugin(plugin)
		if err == nil {
			gs.UpdatePlugin(plugin)
		} else {
			err = gs.InsertPlugin(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dbPlugin, err = gs.GetPlugin(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		plugin.Id = dbPlugin.Id

		err = repo.UploadPlugin(user, plugin, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

	}
}

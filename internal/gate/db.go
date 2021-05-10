package gate

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate/client"
)

func usersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	client := client.NewUserClient("", "")

	switch req.Method {
	case http.MethodGet:
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

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

	client := client.NewPluginClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("name")
		id := r.FormValue("id")
		page := r.FormValue("page")
		count := r.FormValue("count")
		search := r.FormValue("search")

		if pluginName != "" || id != "" {

			req := &api.Plugin{
				Id:   id,
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
				Page:   int32(convPage),
				Count:  int32(convCount),
				Search: search,
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
	client := client.NewReadmeClient("", "")
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

		dbpl, err := gs.GetPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		rdme, err := client.Get(dbpl)
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

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &api.Plugin{
			Id:   r.FormValue("plugin_id"),
			Name: r.FormValue("plugin_name"),
		}

		dbPl, err := gs.GetPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		readme := &api.Readme{
			Plugin: dbPl,
			Text:   r.FormValue("text"),
		}

		err = client.Insert(readme)
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

func sessionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	client := client.NewSessionsClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := r.FormValue("id")
		userId := r.FormValue("userId")

		req := &api.Session{
			Id:     id,
			UserId: userId,
		}
		ses, err := client.Get(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		asJSON, err := json.Marshal(ses)
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
		pl := &api.Session{}

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
	case http.MethodDelete:
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Session{}

		err = json.Unmarshal(bs, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Delete(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

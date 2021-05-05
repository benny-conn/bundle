package gate

import (
	"encoding/json"
	"fmt"
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
		version := r.FormValue("version")

		req := &api.Plugin{
			Name: pluginName,
		}

		dbPl, err := gs.GetPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}

		if version != "latest" {
			dbPl.Version = version
		}
		pl, err := repo.DownloadPlugin(dbPl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(pl)
		return

	case http.MethodPost:

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := &api.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
		plugin := &api.Plugin{
			Name:    r.FormValue("name"),
			Version: r.FormValue("version"),
		}

		dbUser, err := gs.GetUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		plugin.Author = dbUser

		_, err = gs.GetPlugin(plugin)
		if err == nil {
			gs.UpdatePlugin(plugin)
		} else {
			err = gs.InsertPlugin(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		dbPlugin, err := gs.GetPlugin(plugin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("plugin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = repo.UploadPlugin(dbUser, dbPlugin, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

	}
}

func repoThumbnailsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	repo := repo.NewRepoService("", "")
	gs := NewGateService("", "")

	switch r.Method {

	case http.MethodPost:

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := &api.User{
			Id: r.FormValue("user"),
		}
		plugin := &api.Plugin{
			Id: r.FormValue("plugin"),
		}

		dbUser, err := gs.GetUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, h, err := r.FormFile("thumbnail")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if h.Size > (5 << 20) {
			http.Error(w, "file too large", http.StatusBadRequest)
			return
		}

		if plugin.Id == "" {
			fmt.Println("HMM")
			err = repo.UploadThumbnail(dbUser, nil, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		} else {
			// expected end of json input on thumbnail upload why??
			dbPlugin, err := gs.GetPlugin(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = repo.UploadThumbnail(dbUser, dbPlugin, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			dbPlugin.Thumbnail = fmt.Sprintf("https://bundle-repository.s3-us-east-1.amazonaws.com/%s/%s/THUMBNAIL.webp", plugin.Author, plugin.Id)

			err = gs.UpdatePlugin(dbPlugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}

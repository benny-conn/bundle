package gate

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/repo"
)

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

		if version != "latest" && version != "" {
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
			Name:        r.FormValue("name"),
			Version:     r.FormValue("version"),
			Description: r.FormValue("description"),
		}

		if cat, err := strconv.Atoi(r.FormValue("category")); err == nil {
			plugin.Category = api.Category(cat)
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
			err = repo.UploadThumbnail(dbUser, nil, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		} else {

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

			dbPlugin.Thumbnail = fmt.Sprintf("https://bundle-repository.s3.amazonaws.com/%s/%s/THUMBNAIL.webp", dbPlugin.Author.Id, dbPlugin.Id)
			err = gs.UpdatePlugin(dbPlugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}

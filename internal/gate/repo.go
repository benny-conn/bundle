package gate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate/grpc"
	"github.com/bennycio/bundle/internal/repo"
)

func repoPluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	repo := repo.NewRepoService("", "")
	dbcl := grpc.NewPluginClient("", "")
	uscl := grpc.NewUserClient("", "")

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

		dbPl, err := dbcl.Get(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if dbPl.Premium != nil {
			if dbPl.Premium.Price > 0 {
				uclient := grpc.NewUserClient("", "")
				user := r.FormValue("user")
				u := &api.User{}
				err = json.Unmarshal([]byte(user), u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				dbUser, err := uclient.Get(u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if dbUser.Purchases == nil {
					http.Error(w, "user does not own premium plugin", http.StatusUnauthorized)
					return
				}
				can := false
				for _, v := range dbUser.Purchases {
					if v.ObjectId == dbPl.Id && v.Complete {
						can = true
					}
				}
				if !can {
					http.Error(w, "user does not own premium plugin", http.StatusUnauthorized)
					return
				}
			}
		}

		if version != "latest" && version != "" {
			dbPl.Version = version
		}
		pl, err := repo.DownloadPlugin(dbPl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
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

		dbUser, err := uscl.Get(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		plugin.Author = dbUser

		dbPlIni, err := dbcl.Get(plugin)

		if err == nil {
			if dbUser.Id != dbPlIni.Author.Id {
				http.Error(w, "cannot update another author's plugin", http.StatusUnauthorized)
				return
			} else {
				err = dbcl.Update(plugin)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		} else {
			err = dbcl.Insert(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		dbPlugin, err := dbcl.Get(plugin)
		if err != nil {
			fmt.Println(err.Error())
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
			fmt.Println(err.Error())
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

package web

import (
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/russross/blackfriday/v2"
)

const perPageCount = 15

func pluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {

	data := TemplateData{}
	user, err := getProfFromCookie(req)

	if err == nil {
		data.Profile = user
	}

	switch req.Method {

	case http.MethodGet:
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		gs := gate.NewGateService("", "")

		pluginName := req.FormValue("plugin")

		if pluginName == "" {
			page := req.FormValue("page")
			search := req.FormValue("search")

			if page == "" {
				page = "1"
			}

			pageNumber, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			data.Page = pageNumber

			math := func(i int, p int, op string) int {
				result := 0
				switch op {
				case "add":
					result = i + p
				case "sub":
					result = i - p
				case "mul":
					result = i * p
				case "div":
					result = i / p
				}
				return result
			}
			data.Math = math

			plugins, err := gs.PaginatePlugins(pageNumber, perPageCount, search)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			data.Plugins = plugins

		} else {
			req := &api.Plugin{
				Name: pluginName,
			}
			plugin, err := gs.GetPlugin(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			readme, err := gs.GetReadme(req)

			if err == nil {
				output := blackfriday.Run([]byte(readme.Text))
				data.Readme = string(output)
			}

			data.Plugin = plugin
		}
	case http.MethodPost:

		err = req.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := req.FormValue("user")
		plugin := req.FormValue("plugin")

		thumbnail, h, err := req.FormFile("thumbnail")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if h.Size > (1 << 20) {
			http.Error(w, "file too large", http.StatusBadRequest)
			return
		}

		if user == "" {
			http.Error(w, "no user specified", http.StatusBadRequest)
			return
		}
		if plugin == "" {
			http.Error(w, "no user specified", http.StatusBadRequest)
			return
		}

		gs := gate.NewGateService("", "")

		u := &api.User{
			Id: user,
		}
		p := &api.Plugin{
			Id: plugin,
		}

		err = gs.UploadThumbnail(u, p, thumbnail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, req, req.Header.Get("Referer"), http.StatusFound)
		return
	}

	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

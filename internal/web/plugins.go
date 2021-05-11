package web

import (
	"net/http"
	"strconv"
	"time"

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
			sort := req.FormValue("sort")
			category := req.FormValue("category")

			if page == "" {
				page = "1"
			}
			if sort == "" {
				sort = "0"
			}
			if category == "" {
				category = "0"
			}

			pageNumber, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			sortNumber, err := strconv.Atoi(sort)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			categoryNum, err := strconv.Atoi(category)
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

			date := func(i int64) string {
				d := time.Unix(i, 0)
				return d.Format("Mon Jan 2 15:04:05 MST")
			}
			data.Date = date

			req := &api.PaginatePluginsRequest{
				Count:    10,
				Search:   search,
				Page:     int32(pageNumber),
				Sort:     api.Sort(sortNumber),
				Category: api.Category(categoryNum),
			}

			plugins, err := gs.PaginatePlugins(req)
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

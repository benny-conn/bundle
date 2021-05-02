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

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUserFromCookie(req)

	data := TemplateData{}
	if err == nil {
		data.User = user
	}

	gs := gate.NewGateService("", "")

	pluginName := req.FormValue("plugin")

	if pluginName == "" {
		page := req.FormValue("page")

		if page == "" {
			page = "1"
		}

		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		plugins, err := gs.PaginatePlugins(pageNumber, perPageCount)
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
	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		panic(err)
	}
}

package web

import (
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/wrapper"
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

		plugins, err := wrapper.PaginatePluginsApi(pageNumber, perPageCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data.Plugins = plugins

		err = tpl.ExecuteTemplate(w, "plugins", data)
		if err != nil {
			panic(err)
		}

	} else {
		req := &api.Plugin{
			Name: pluginName,
		}
		plugin, err := wrapper.GetPluginApi(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		readme, err := wrapper.DownloadReadmeApi(req)

		if err == nil {
			output := blackfriday.Run(readme)
			plugin.Readme = string(output)
		}

		data.Plugin = plugin

		err = tpl.ExecuteTemplate(w, "plugins", data)
		if err != nil {
			panic(err)
		}
	}
}

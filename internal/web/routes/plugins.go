package routes

import (
	"net/http"
	"strconv"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
)

func PluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _ := getProfileFromCookie(req)

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

		plugins, err := wrapper.PaginatePluginsApi(pageNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data := bundle.TemplateData{
			Plugins: plugins,
			User:    *user,
		}

		err = tpl.ExecuteTemplate(w, "plugin", data)
		if err != nil {
			panic(err)
		}
		return

	}

	plugin, err := wrapper.GetPluginApi(pluginName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// DO THIS
	// md, err := pkg.GetPluginData(opts)
	// if err == nil {
	// 	output := blackfriday.Run(md.Readme)
	// }

	data := bundle.TemplateData{
		User:   *user,
		Plugin: *plugin,
	}

	err = tpl.ExecuteTemplate(w, "plugin", data)
	if err != nil {
		panic(err)
	}
}

package routes

import (
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/pkg"
	"github.com/russross/blackfriday/v2"
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

		plugins, err := pkg.PaginatePlugins(pageNumber)
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

	plugin, err := pkg.GetPlugin(pluginName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	opts := &api.GetPluginDataRequest{
		Name:       plugin.Name,
		Version:    plugin.Version,
		WithReadme: true,
	}

	// make readme not bytes
	md, err := pkg.GetPluginData(opts)
	if err == nil {
		output := blackfriday.Run(md.Readme)
		plugin.Readme = string(output)
	}

	data := bundle.TemplateData{
		User:   *user,
		Plugin: *plugin,
	}

	err = tpl.ExecuteTemplate(w, "plugin", data)
	if err != nil {
		panic(err)
	}
}

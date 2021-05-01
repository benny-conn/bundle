package web

import (
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/wrapper"
)

const perPageCount = 15

func pluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _ := getUserFromCookie(req)

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

		data := TemplateData{
			Plugins: plugins,
			User:    user,
		}

		err = tpl.ExecuteTemplate(w, "plugins", data)
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

	// readme, err := wrapper.DownloadReadmeApi(pluginName)

	// if err == nil {
	// 	output := blackfriday.Run(readme)
	// 	pluginInfo.Readme = string(output)
	// }

	data := TemplateData{
		User:   user,
		Plugin: plugin,
	}

	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		panic(err)
	}
}

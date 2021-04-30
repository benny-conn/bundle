package web

import (
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/wrapper"
	"github.com/russross/blackfriday/v2"
)

type PluginInfo struct {
	Name        string
	Author      string
	Version     string
	Description string
	Readme      string
	Thumbnail   []byte
}

func pluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {

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

		infos := pluginsToInfos(plugins)

		data := TemplateData{
			Plugins: infos,
			Profile: user,
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

	pluginInfo := pluginToInfo(plugin)

	readme, err := wrapper.DownloadReadmeApi(pluginName)

	if err == nil {
		output := blackfriday.Run(readme)
		pluginInfo.Readme = string(output)
	}

	data := TemplateData{
		Profile: user,
		Plugin:  pluginInfo,
	}

	err = tpl.ExecuteTemplate(w, "plugins", data)
	if err != nil {
		panic(err)
	}
}

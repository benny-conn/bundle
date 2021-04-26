package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/internal/storage"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

const ReqFileType = bundle.RequiredFileType

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func RootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	profile, _ := getProfileFromCookie(r)
	data := bundle.TemplateData{
		Profile: profile,
	}

	err := tpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		panic(err)
	}

}
func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, _ := getProfileFromCookie(r)
	err := tpl.ExecuteTemplate(w, "register", user)
	if err != nil {
		bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func BundleHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")

		plugin, err := storage.GetPlugin(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if version != "latest" {
			plugin.Version = version
		}

		bs, err := storage.DownloadFromRepo(plugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		bundle.WriteResponse(w, string(bs), http.StatusOK)
	}

	if r.Method == http.MethodPost {

		plugin := r.Header.Get("Plugin")

		reqPlugin := &bundle.Plugin{}

		err := json.Unmarshal([]byte(plugin), reqPlugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		finalPlugin, err := validateAndReturnPlugin(*reqPlugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		uploadLocation, err := storage.UploadToRepo(finalPlugin, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to", uploadLocation)
	}

}

func UserHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		newUser := bundle.User{
			Username: req.FormValue("username"),
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		}
		err := storage.InsertUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

}

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("plugin")

		plugin, err := storage.GetPlugin(pluginName)
		if err != nil {
			panic(err)
		}

		asJSON, err := json.Marshal(plugin)
		if err != nil {
			panic(err)
		}

		bundle.WriteResponse(w, string(asJSON), http.StatusOK)
	}
}

func LoginHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	if req.Method == http.MethodPost {

		req.ParseForm()
		user := bundle.User{
			Username: req.FormValue("username"),
			Password: req.FormValue("password"),
		}

		asJSON, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userJSON := string(asJSON)
		validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := storage.GetUser(*validatedUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(validatedUser.Password))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		profile := bundle.Profile{
			Username: dbUser.Username,
			Email:    dbUser.Email,
			Scopes:   dbUser.Scopes,
		}

		token, err := auth.NewAuthToken(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tokenCookie := auth.NewAccessCookie(token)
		http.SetCookie(w, tokenCookie)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	}

	if req.Method == http.MethodGet {
		user, _ := getProfileFromCookie(req)
		err := tpl.ExecuteTemplate(w, "login", user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LogoutHandlerFunc(w http.ResponseWriter, req *http.Request) {

	accessCookie, err := req.Cookie("access_token")
	if err == nil {
		accessCookie.MaxAge = -1
	}

	http.SetCookie(w, accessCookie)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func PluginHandlerFunc(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, _ := getProfileFromCookie(req)

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

		plugins, err := storage.PaginatePlugins(pageNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data := bundle.TemplateData{
			Plugins: plugins,
			Profile: profile,
		}

		err = tpl.ExecuteTemplate(w, "plugin", data)
		if err != nil {
			panic(err)
		}
		return

	}

	plugin, err := storage.GetPlugin(pluginName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	pluginReadme := bundle.Plugin{
		Name:    plugin.Name,
		Author:  plugin.Author,
		Version: "README",
	}

	md, err := storage.DownloadFromRepo(pluginReadme)
	if err == nil {
		output := blackfriday.Run(md)
		plugin.Readme = string(output)
	}

	data := bundle.TemplateData{
		Profile: profile,
		Plugin:  plugin,
	}

	err = tpl.ExecuteTemplate(w, "plugin", data)
	if err != nil {
		panic(err)
	}
}

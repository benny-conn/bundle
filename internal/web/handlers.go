package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

		version := r.Header.Get("Plugin-Version")
		userJSON := r.Header.Get("User")
		pluginName := r.Header.Get("Plugin-Name")

		isReadme := strings.Contains(version, "README")

		validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := storage.GetUser(*validatedUser)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqPlugin := bundle.Plugin{
			Name:    pluginName,
			Author:  dbUser.Username,
			Version: version,
		}

		dbPlugin, err := storage.GetPlugin(pluginName)

		if err == nil {
			isUserPluginAuthor := strings.EqualFold(dbPlugin.Author, validatedUser.Username)

			if isUserPluginAuthor {
				if isReadme {
					dbPlugin.Version = version
				} else {
					err = storage.UpdatePlugin(pluginName, reqPlugin)
					if err != nil {
						bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				err = errors.New("you are not permitted to edit this plugin")
				bundle.WriteResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
		} else {
			if isReadme {
				err = errors.New("no plugin to attach readme to")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = storage.InsertPlugin(reqPlugin)
			dbPlugin = reqPlugin
			if err != nil {
				bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		uploadLocation, err := storage.UploadToRepo(dbPlugin, r.Body)
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	output := blackfriday.Run(md)

	user, _ := getProfileFromCookie(req)

	data := bundle.TemplateData{
		Profile:  user,
		Plugin:   plugin,
		Markdown: string(output),
	}

	err = tpl.ExecuteTemplate(w, "plugin", data)
	if err != nil {
		panic(err)
	}
}

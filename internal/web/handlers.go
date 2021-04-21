package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

const ReqFileType = bundle.RequiredFileType

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*.gohtml"))
}

func RootHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func SignupHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
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

		plugin, err := bundle.GetPlugin(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if version != "latest" {
			plugin.Version = version
		}

		bs, err := bundle.DownloadPlugin(*plugin)
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

		validatedUser, err := bundle.ValidateAndReturnUser(userJSON)

		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := bundle.GetUser(*validatedUser)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		reqPlugin := &bundle.Plugin{
			Plugin:  pluginName,
			User:    dbUser.Username,
			Version: version,
		}

		decodedPluginResult, err := bundle.GetPlugin(pluginName)

		if err == nil {
			isUserPluginAuthor := strings.EqualFold(decodedPluginResult.User, validatedUser.Username)

			if isUserPluginAuthor {
				err = bundle.UpdatePlugin(pluginName, *reqPlugin)
				if err != nil {
					bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				err = errors.New("you are not permitted to edit this plugin")
				bundle.WriteResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
		} else {
			err = bundle.InsertPlugin(*reqPlugin)
			if err != nil {
				bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			decodedPluginResult.Plugin = pluginName
			decodedPluginResult.User = dbUser.Username
			decodedPluginResult.Version = version
		}

		uploadLocation, err := bundle.UploadPlugin(*decodedPluginResult, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to", uploadLocation)
	}

}

func UserHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == http.MethodPost {

		r.ParseForm()
		newUser := bundle.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		err := bundle.InsertUser(newUser)
		if err != nil {
			bundle.WriteResponse(w, "error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		bundle.WriteResponse(w, "", http.StatusOK)
	}

}

func PluginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("plugin")

		plugin, err := bundle.GetPlugin(pluginName)
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

	req.ParseForm()
	user := bundle.User{
		Username: req.FormValue("username"),
		Email:    req.FormValue("email"),
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

	dbUser, err := bundle.GetUser(*validatedUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(validatedUser.Password))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// TODO make this part of database interaction
	role := "user"

	token, err := auth.GetAuthToken(role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	c := &http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(w, c)
}

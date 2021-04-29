package routes

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/wrapper"
)

func ProfileHandlerFunc(w http.ResponseWriter, req *http.Request) {

	user, err := getProfileFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println(user)
	if req.Method == http.MethodPost {

		req.ParseForm()

		newUsername := req.FormValue("username")
		newTag := req.FormValue("tag")

		updatedUser := &api.User{
			Username: newUsername,
			Tag:      newTag,
		}

		err = wrapper.UpdateUserApi(user.Username, updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbUpdatedUser, _ := wrapper.GetUserApi(updatedUser.Username, updatedUser.Email)

		fmt.Println(dbUpdatedUser)

		token, _ := auth.NewAuthToken(dbUpdatedUser)
		c := auth.NewAccessCookie(token)
		http.SetCookie(w, c)

	}

	data := bundle.TemplateData{
		Profile: user,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

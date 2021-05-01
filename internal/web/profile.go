package web

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/api"

	"github.com/bennycio/bundle/wrapper"
)

func profileHandlerFunc(w http.ResponseWriter, req *http.Request) {

	user, err := getUserFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		req.ParseForm()

		newUsername := req.FormValue("username")
		newTag := req.FormValue("tag")

		updatedUser := &api.User{
			Username: newUsername,
			Tag:      newTag,
		}

		err = wrapper.UpdateUserApi(updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbUpdatedUser, _ := wrapper.GetUserApi(updatedUser.Username, updatedUser.Email)

		fmt.Println(dbUpdatedUser)

		token, _ := newAuthToken(dbUpdatedUser)
		c := newAccessCookie(token)
		http.SetCookie(w, c)

	}

	data := TemplateData{
		User: user,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

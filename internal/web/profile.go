package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
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
			Id:       user.Id,
			Username: newUsername,
			Tag:      newTag,
		}

		gs := gate.NewGateService("", "")
		err = gs.UpdateUser(updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// NOW PROPOGATE THOSE CHANGE BOYO!!!

		dbUpdatedUser, _ := gs.GetUser(updatedUser)

		token, _ := newAuthToken(dbUpdatedUser)
		c := newAccessCookie(token)
		http.SetCookie(w, c)
		user = dbUpdatedUser
	}

	data := TemplateData{
		User: user,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

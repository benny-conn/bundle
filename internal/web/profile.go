package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

func profileHandlerFunc(w http.ResponseWriter, req *http.Request) {

	pro, err := getProfFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		req.ParseForm()

		newUsername := req.FormValue("username")
		newTag := req.FormValue("tag")

		updatedUser := &api.User{
			Id:       pro.Id,
			Username: newUsername,
			Tag:      newTag,
		}

		gs := gate.NewGateService("", "")
		err = gs.UpdateUser(updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbUpdatedUser, _ := gs.GetUser(updatedUser)

		updatedProfile := userToProfile(dbUpdatedUser)

		token, _ := newSession(updatedProfile)
		c := newAccessCookie(token.Id)
		http.SetCookie(w, c)
		pro = updatedProfile
	}

	data := TemplateData{
		Profile: pro,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

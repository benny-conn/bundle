package web

import (
	"net/http"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

func logoutHandlerFunc(w http.ResponseWriter, req *http.Request) {

	accessCookie, err := req.Cookie("access_token")
	if err == nil {
		accessCookie.MaxAge = -1
		gs := gate.NewGateService("", "")

		gs.DeleteBundle(&api.Bundle{Id: accessCookie.Value})
	}

	http.SetCookie(w, accessCookie)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

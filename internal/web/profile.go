package web

import (
	"net/http"
	"strconv"

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

		ftpUser := req.FormValue("ftp-username")
		ftpPass := req.FormValue("ftp-password")
		ftpPort := req.FormValue("ftp-port")
		ftpHost := req.FormValue("ftp-host")

		port, err := strconv.Atoi(ftpPort)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Bundle{
			UserId:  pro.Id,
			FtpUser: ftpUser,
			FtpPass: ftpPass,
			FtpPort: int32(port),
			FtpHost: ftpHost,
		}

		gs := gate.NewGateService("", "")

		_, err = gs.GetBundle(req)

		if err == nil {
			err = gs.UpdateBundle(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {

			err = gs.InsertBundle(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	data := TemplateData{
		Profile: pro,
	}

	err = tpl.ExecuteTemplate(w, "profile", data)
	if err != nil {
		panic(err)
	}

}

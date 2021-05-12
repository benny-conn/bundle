package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/jlaffaye/ftp"

	"golang.org/x/crypto/bcrypt"
)

func ftpHandlerFunc(w http.ResponseWriter, req *http.Request) {

	pro, err := getProfFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		req.ParseForm()
		pass := req.FormValue("ftp-password")

		gs := gate.NewGateService("", "")
		dbBundle, err := gs.GetBundle(&api.Bundle{UserId: pro.Id})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbBundle.FtpPass), []byte(pass))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		c, err := ftp.Dial(fmt.Sprintf("%s:%v", dbBundle.FtpHost, dbBundle.FtpPort), ftp.DialWithTimeout(10*time.Second))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = c.Login(dbBundle.FtpUser, dbBundle.FtpPass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := c.Retr("bukkit.yml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bf := bytes.Buffer{}

		_, err = io.Copy(w, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := c.Quit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bf.Bytes())

	}

	if req.Method == http.MethodGet {

		referer := req.Referer()

		err := tpl.ExecuteTemplate(w, "ftp", referer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

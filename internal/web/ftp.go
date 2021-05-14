package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

	gs := gate.NewGateService("", "")

	data := TemplateData{
		Profile: pro,
	}

	if req.Method == http.MethodPost {

		req.ParseForm()
		pass := req.FormValue("ftp-password")
		pls := req.FormValue("plugins")

		plsSplit := strings.Split(pls, ",")

		data.Profile.Bundles = plsSplit

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

		dbBundle.Plugins = plsSplit

		err = gs.UpdateBundle(dbBundle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c, err := ftp.Dial(fmt.Sprintf("%s:%v", dbBundle.FtpHost, dbBundle.FtpPort), ftp.DialWithTimeout(10*time.Second))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = c.Login(dbBundle.FtpUser, pass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = c.ChangeDir("plugins")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wg := &sync.WaitGroup{}
		wg.Add(len(dbBundle.Plugins))
		for _, v := range dbBundle.Plugins {
			go func(pl string) {
				defer wg.Done()
				c.Delete(pl + ".jar")
				bs, err := gs.DownloadPlugin(&api.Plugin{Name: pl})
				if err != nil {
					return
				}
				buf := bytes.NewBuffer(bs)
				err = c.Stor(pl+".jar", buf)
				fmt.Println(err.Error())
			}(v)
		}

		wg.Wait()
		if err := c.Quit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	dbBundle, err := gs.GetBundle(&api.Bundle{UserId: pro.Id})
	if err == nil {
		data.Profile.Bundles = dbBundle.Plugins
	}

	err = tpl.ExecuteTemplate(w, "ftp", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

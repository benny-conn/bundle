package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/jlaffaye/ftp"
)

func ftpHandlerFunc(w http.ResponseWriter, req *http.Request) {

	pro, err := getProfFromCookie(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	gs := gate.NewGateService("", "")

	data := templateData{
		Profile: pro,
	}

	if req.Method == http.MethodPost {

		req.ParseForm()
		pls := req.FormValue("plugins")

		fmt.Println(pls)

		plsSplit := strings.Split(pls, ",")

		ftpUser := req.FormValue("ftp-username")
		ftpPass := req.FormValue("ftp-password")
		ftpPort := req.FormValue("ftp-port")
		ftpHost := req.FormValue("ftp-host")
		save := req.FormValue("save")

		port, err := strconv.Atoi(ftpPort)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sBool := false
		if save == "on" {
			sBool = true
		}

		req := &api.Bundle{
			UserId:  pro.Id,
			FtpUser: ftpUser,
			FtpPort: int32(port),
			FtpHost: ftpHost,
			Plugins: plsSplit,
		}

		if sBool {
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

		go func() {
			c, err := ftp.Dial(fmt.Sprintf("%s:%v", req.FtpHost, req.FtpPort), ftp.DialWithTimeout(10*time.Second))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = c.Login(req.FtpUser, ftpPass)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = c.ChangeDir("plugins")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var wg *sync.WaitGroup
			wg.Add(len(req.Plugins))
			for _, v := range req.Plugins {
				go func(pl string) {
					defer wg.Done()
					c.Delete(pl + ".jar")
					bs, err := gs.DownloadPlugin(&api.Plugin{Name: pl})
					if err != nil {
						return
					}
					buf := bytes.NewBuffer(bs)
					err = c.Stor(pl+".jar", buf)
					if err != nil {
						fmt.Println(err.Error())
					}
				}(v)
			}

			wg.Wait()
			if err := c.Quit(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

	}

	dbBundle, err := gs.GetBundle(&api.Bundle{UserId: pro.Id})
	if err == nil {
		data.Bundle = dbBundle
	}

	err = tpl.ExecuteTemplate(w, "ftp", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

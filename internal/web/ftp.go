package web

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

		plsSplit := strings.Split(pls, ",")

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
			FtpPort: int32(port),
			FtpHost: ftpHost,
			Plugins: plsSplit,
		}

		err = gs.UpdateBundle(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		go func(req *api.Bundle) {
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
			for _, v := range req.Plugins {
				bs, err := gs.DownloadPlugin(&api.Plugin{Name: v})
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return
				}
				rd := &bytes.Buffer{}
				n, err := rd.Write(bs)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return
				}
				fmt.Println(n)
				err = c.Stor(v+".jar", rd)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return
				}
				fmt.Println(n)
			}

			if err := c.Quit(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}(req)

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

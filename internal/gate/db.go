package gate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate/grpc"
)

func usersHandlerFunc(w http.ResponseWriter, req *http.Request) {

	client := grpc.NewUserClient("", "")

	switch req.Method {
	case http.MethodGet:
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		userName := req.FormValue("username")
		email := req.FormValue("email")
		id := req.FormValue("id")

		r := &api.User{
			Id:       id,
			Username: userName,
			Email:    email,
		}
		user, err := client.Get(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bs, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bs)
	case http.MethodPost:

		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser := &api.User{}
		err = json.Unmarshal(bs.Bytes(), newUser)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = client.Insert(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPatch:
		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser := &api.User{}
		err = json.Unmarshal(bs.Bytes(), newUser)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = client.Update(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func pluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	client := grpc.NewPluginClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("name")
		id := r.FormValue("id")
		page := r.FormValue("page")
		count := r.FormValue("count")
		search := r.FormValue("search")

		if pluginName != "" || id != "" {

			req := &api.Plugin{
				Id:   id,
				Name: pluginName,
			}

			plugin, err := client.Get(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			asJSON, err := json.Marshal(plugin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		} else if page != "" {
			convPage, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			convCount, err := strconv.Atoi(count)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			req := &api.PaginatePluginsRequest{
				Page:   int32(convPage),
				Count:  int32(convCount),
				Search: search,
			}
			plugins, err := client.Paginate(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			asJSON, err := json.Marshal(plugins)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		}
	case http.MethodPost:
		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pl := &api.Plugin{}

		err = json.Unmarshal(bs.Bytes(), pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Insert(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Plugin{}

		err = json.Unmarshal(bs.Bytes(), req)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Update(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func readmesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	client := grpc.NewReadmeClient("", "")
	gs := NewGateService("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pluginName := r.FormValue("name")
		id := r.FormValue("id")

		req := &api.Plugin{
			Id:   id,
			Name: pluginName,
		}

		dbpl, err := gs.GetPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		rdme, err := client.Get(dbpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		asJSON, err := json.Marshal(rdme)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		internal.WriteResponse(w, string(asJSON), http.StatusOK)
		return

	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &api.Plugin{
			Id:   r.FormValue("plugin_id"),
			Name: r.FormValue("plugin_name"),
		}

		dbPl, err := gs.GetPlugin(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		readme := &api.Readme{
			Plugin: dbPl,
			Text:   r.FormValue("text"),
		}

		_, err = gs.GetReadme(dbPl)
		if err == nil {
			err = client.Update(readme)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = client.Insert(readme)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPatch:

		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Readme{}

		err = json.Unmarshal(bs.Bytes(), req)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Update(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func sessionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	client := grpc.NewSessionsClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := r.FormValue("id")
		userId := r.FormValue("userId")

		req := &api.Session{
			Id:     id,
			UserId: userId,
		}
		ses, err := client.Get(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		asJSON, err := json.Marshal(ses)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		internal.WriteResponse(w, string(asJSON), http.StatusOK)
		return

	case http.MethodPost:

		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pl := &api.Session{}
		if err = json.Unmarshal(bs.Bytes(), pl); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := client.Insert(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		internal.WriteResponse(w, string(r), http.StatusOK)
	case http.MethodDelete:
		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Session{}

		err = json.Unmarshal(bs.Bytes(), req)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Delete(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func changelogHandlerFunc(w http.ResponseWriter, r *http.Request) {
	client := grpc.NewChangelogsClient("", "")

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := r.FormValue("id")
		pluginId := r.FormValue("pluginId")
		version := r.FormValue("version")

		if version == "" {
			req := &api.Changelog{
				PluginId: pluginId,
			}
			ses, err := client.GetAll(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			asJSON, err := json.Marshal(ses)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		} else {
			req := &api.Changelog{
				Id:       id,
				PluginId: pluginId,
				Version:  version,
			}
			ses, err := client.Get(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			asJSON, err := json.Marshal(ses)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		}

	case http.MethodPost:

		bs := &bytes.Buffer{}
		_, err := io.Copy(bs, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bu := &api.Changelog{}

		err = json.Unmarshal(bs.Bytes(), bu)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = client.Insert(bu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"google.golang.org/grpc"
)

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		pluginName := r.FormValue("name")
		page := r.FormValue("page")

		if pluginName != "" {
			plugin, err := getPlugin(pluginName)
			if err != nil {
				panic(err)
			}

			asJSON, err := json.Marshal(plugin)
			if err != nil {
				panic(err)
			}

			internal.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		} else if page != "" {
			convPage, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			plugins, err := paginatePlugins(convPage)
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
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pl := &api.Plugin{}

		err = json.Unmarshal(bs, pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = insertPlugin(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pl := &api.UpdatePluginRequest{}

		err = json.Unmarshal(bs, pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = updatePlugin(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func getPlugin(name string) (*api.Plugin, error) {

	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	req := &api.GetPluginRequest{Name: name}
	pl, err := client.GetPlugin(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func updatePlugin(req *api.UpdatePluginRequest) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)

	_, err = client.UpdatePlugin(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func insertPlugin(plugin *api.Plugin) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	_, err = client.InsertPlugin(context.Background(), plugin)
	if err != nil {
		return err
	}
	return nil
}

func paginatePlugins(page int) ([]*api.Plugin, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	req := &api.PaginatePluginsRequest{}
	results, err := client.PaginatePlugins(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return results.Plugins, nil
}

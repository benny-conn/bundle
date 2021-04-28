package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	pb "github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"google.golang.org/grpc"
)

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
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

			bundle.WriteResponse(w, string(asJSON), http.StatusOK)
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
			bundle.WriteResponse(w, string(asJSON), http.StatusOK)
			return
		}
	}
}

func getPlugin(name string) (*pb.Plugin, error) {

	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewPluginsServiceClient(conn)
	req := &pb.GetPluginRequest{Name: name}
	pl, err := client.GetPlugin(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func updatePlugin(name string, updatedPlugin *pb.Plugin) error {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewPluginsServiceClient(conn)
	req := &pb.UpdatePluginRequest{Name: name, UpdatedPlugin: updatedPlugin}
	_, err = client.UpdatePlugin(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func insertPlugin(plugin *pb.Plugin) error {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewPluginsServiceClient(conn)
	_, err = client.InsertPlugin(context.Background(), plugin)
	if err != nil {
		return err
	}
	return nil
}

func paginatePlugins(page int) ([]*pb.Plugin, error) {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewPluginsServiceClient(conn)
	req := &pb.PaginatePluginsRequest{}
	results, err := client.PaginatePlugins(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return results.Plugins, nil
}

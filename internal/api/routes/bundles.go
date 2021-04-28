package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"google.golang.org/grpc"
)

func BundleHandlerFunc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")
		withPlugin, err := strconv.ParseBool(r.FormValue("withPlugin"))
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		withReadme, err := strconv.ParseBool(r.FormValue("withReadme"))
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		withThumbnail, err := strconv.ParseBool(r.FormValue("withThumbnail"))
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		opts := &api.GetPluginDataRequest{
			Name:          name,
			Version:       version,
			WithPlugin:    withPlugin,
			WithReadme:    withReadme,
			WithThumbnail: withThumbnail,
		}

		plugin, err := getPlugin(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if version != "latest" {
			plugin.Version = version
		}

		pl, err := downloadFromRepo(opts)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		bs, err := json.Marshal(pl)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bundle.WriteResponse(w, string(bs), http.StatusOK)
	}

	if r.Method == http.MethodPost {

		dataType := r.Header.Get("Data-Type")

		pluginJSON := r.Header.Get("Resource")

		reqData := &api.InsertPluginDataRequest{}

		plugin := &api.Plugin{}
		json.Unmarshal([]byte(pluginJSON), plugin)

		bs, err := io.ReadAll(r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch dataType {
		case api.InsertPluginDataRequest_PLUGIN.String():
			reqData.DataType = api.InsertPluginDataRequest_PLUGIN
			reqData.Plugin.PluginData = bs
		case api.InsertPluginDataRequest_README.String():
			reqData.DataType = api.InsertPluginDataRequest_README
			reqData.Plugin.Readme = bs
		case api.InsertPluginDataRequest_THUMBNAIL.String():
			reqData.DataType = api.InsertPluginDataRequest_THUMBNAIL
			reqData.Plugin.Thumbnail = bs
		default:
			reqData.DataType = api.InsertPluginDataRequest_PLUGIN
			reqData.Plugin.PluginData = bs
		}

		err = updateOrInsertPlugin(reqData)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = uploadToRepo(reqData)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded " + reqData.Plugin.Name)
	}

}

func updateOrInsertPlugin(req *api.InsertPluginDataRequest) error {
	dbPlugin, err := getPlugin(req.Plugin.Name)

	if err == nil {
		err = updatePlugin(dbPlugin.Name, req.Plugin)
		if err != nil {
			return err
		}
	} else {
		if req.DataType.String() != api.InsertPluginDataRequest_PLUGIN.String() {
			return errors.New("no plugin to attach file to")
		}
		err = insertPlugin(req.Plugin)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func downloadFromRepo(req *api.GetPluginDataRequest) (*api.Plugin, error) {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	pl, err := client.GetPluginData(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func uploadToRepo(req *api.InsertPluginDataRequest) error {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	_, err = client.InsertPluginData(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

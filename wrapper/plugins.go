package wrapper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

func PaginatePluginsApi(page int) ([]api.Plugin, error) {

	u, err := url.Parse(ApiServerHost)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", string(page))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &[]api.Plugin{}

	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}
	return *result, nil

}

func GetPluginApi(name string) (*api.Plugin, error) {

	u, err := url.Parse(ApiServerHost + "/plugins")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("name", name)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.Plugin{}

	err = json.Unmarshal(bs, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetPlugin(name string) (*api.Plugin, error) {

	conn, err := grpc.Dial(grpcAddress)
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

func UpdatePlugin(name string, updatedPlugin *api.Plugin) error {
	conn, err := grpc.Dial(grpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	req := &api.UpdatePluginRequest{Name: name, UpdatedPlugin: updatedPlugin}
	_, err = client.UpdatePlugin(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func InsertPlugin(plugin *api.Plugin) error {
	conn, err := grpc.Dial(grpcAddress)
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

func PaginatePlugins(page int) ([]*api.Plugin, error) {
	conn, err := grpc.Dial(grpcAddress)
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

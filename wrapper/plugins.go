package wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

func PaginatePluginsApi(page int) ([]api.Plugin, error) {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf(":%v", port)
	u, err := url.Parse(addr + "/api/plugins")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", fmt.Sprint(page))
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
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf(":%v", port)
	u, err := url.Parse(addr + "/api/plugins")
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

	creds, err := GetCert()
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

func UpdatePlugin(name string, updatedPlugin *api.Plugin) error {
	creds, err := GetCert()
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
	req := &api.UpdatePluginRequest{Name: name, UpdatedPlugin: updatedPlugin}
	_, err = client.UpdatePlugin(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func InsertPlugin(plugin *api.Plugin) error {
	creds, err := GetCert()
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

func PaginatePlugins(page int) ([]*api.Plugin, error) {
	creds, err := GetCert()
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

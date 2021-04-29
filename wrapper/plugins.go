package wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
)

func PaginatePluginsApi(page int) ([]*api.Plugin, error) {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v", port)
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

	result := &[]*api.Plugin{}

	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}
	return *result, nil

}

func GetPluginApi(name string) (*api.Plugin, error) {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v", port)
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

func InsertPluginApi(plugin *api.Plugin) error {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v/api/plugins", port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(plugin)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(bs)

	resp, err := http.Post(u.String(), "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func UpdatePluginApi(pluginName string, updatedPlugin *api.Plugin) error {
	port := os.Getenv("API_PORT")
	addr := fmt.Sprintf("http://localhost:%v/api/plugins", port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	up := &api.UpdatePluginRequest{
		Name:          pluginName,
		UpdatedPlugin: updatedPlugin,
	}

	bs, err := json.Marshal(up)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(bs)

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

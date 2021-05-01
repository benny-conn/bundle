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

func PaginatePluginsApi(page int, count int) ([]*api.Plugin, error) {
	port := os.Getenv("API_PORT")
	host := os.Getenv("API_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/plugins", host, port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", fmt.Sprint(page))
	q.Set("count", fmt.Sprint(count))
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

	result := &api.PaginatePluginsResponse{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}
	return result.Plugins, nil

}

func GetPluginApi(plugin *api.Plugin) (*api.Plugin, error) {
	port := os.Getenv("API_PORT")
	host := os.Getenv("API_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/plugins", host, port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("name", plugin.Name)
	q.Set("id", plugin.Id)
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
	host := os.Getenv("API_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/plugins", host, port)
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

func UpdatePluginApi(updatedPlugin *api.Plugin) error {
	port := os.Getenv("API_PORT")
	host := os.Getenv("API_HOST")
	addr := fmt.Sprintf("http://%v:%v/api/plugins", host, port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(updatedPlugin)
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

func UploadReadmeApi(user *api.User, plugin *api.Plugin, data io.Reader) error {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	if port == "" {
		port = "8060"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/readmes", host, port))
	if err != nil {
		return err
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	pluginJSON, err := json.Marshal(plugin)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), data)
	if err != nil {
		return err
	}
	req.Header.Add("User", string(userJSON))
	req.Header.Add("Resource", string(pluginJSON))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func DownloadReadmeApi(plugin *api.Plugin) ([]byte, error) {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	if port == "" {
		port = "8060"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/readmes", host, port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", plugin.Name)
	q.Add("id", plugin.Id)
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
	return bs, nil
}

package wrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
)

func DownloadReadmeApi(pluginName string) ([]byte, error) {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/readmes", host, port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", pluginName)
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

func DownloadPluginApi(pluginName string, version string) ([]byte, error) {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", host, port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", pluginName)
	q.Add("version", version)
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

func UploadPluginApi(user *api.User, pluginName string, version string, data io.Reader) error {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", host, port))
	if err != nil {
		return err
	}
	fmt.Println(u.String())
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	plugin := &api.Plugin{
		Name:    pluginName,
		Version: version,
		Author:  user.Username,
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

func UploadReadmeApi(user *api.User, pluginName string, data io.Reader) error {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/readmes", host, port))
	if err != nil {
		return err
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	plugin := &api.Plugin{
		Name:   pluginName,
		Author: user.Username,
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

package wrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	auth "github.com/bennycio/bundle/internal/api"
)

func DownloadPluginApi(plugin *api.Plugin) ([]byte, error) {

	port := os.Getenv("API_PORT")
	host := os.Getenv("API_HOST")
	if port == "" {
		port = "8020"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/bundles/plugins", host, port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", plugin.Name)
	q.Add("version", plugin.Version)
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

func UploadPluginApi(user *api.User, plugin *api.Plugin, data io.Reader) error {

	port := os.Getenv("API_PORT")
	host := os.Getenv("API_HOST")
	if port == "" {
		port = "8020"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/bundles/plugins", host, port))
	if err != nil {
		return err
	}
	fmt.Println(u.String())
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

func DownloadPluginRepo(plugin *api.Plugin) ([]byte, error) {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	if port == "" {
		port = "8060"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", host, port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", plugin.Name)
	q.Add("version", plugin.Version)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	access, err := auth.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", "Bearer "+access)
	resp, err := http.DefaultClient.Do(req)
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

func UploadPluginRepo(user *api.User, plugin *api.Plugin, data io.Reader) error {
	port := os.Getenv("REPO_PORT")
	host := os.Getenv("REPO_HOST")
	if port == "" {
		port = "8060"
	}
	if host == "" {
		host = "localhost"
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", host, port))
	if err != nil {
		return err
	}
	fmt.Println(u.String())
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
	access, err := auth.GetAccessToken()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+access)
	req.Header.Add("User", string(userJSON))
	req.Header.Add("Resource", string(pluginJSON))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

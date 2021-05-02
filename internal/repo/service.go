package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
)

type repoService interface {
	DownloadPlugin(plugin *api.Plugin) ([]byte, error)
	UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error
}

type repoServiceImpl struct {
	Host string
	Port string
}

func NewRepoService(host string, port string) repoService {
	if host == "" {
		host = os.Getenv("REPO_HOST")
	}
	if port == "" {
		port = os.Getenv("REPO_PORT")
	}
	return &repoServiceImpl{
		Host: host,
		Port: port,
	}
}

func (r *repoServiceImpl) DownloadPlugin(plugin *api.Plugin) ([]byte, error) {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", r.Host, r.Port))
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
	access, err := getAccessToken()
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

func (r *repoServiceImpl) UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/repo/plugins", r.Host, r.Port))
	if err != nil {
		return err
	}
	fmt.Println(u.String())

	fmt.Println("WE UPLOADING")
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	pluginJSON, err := json.Marshal(plugin)
	if err != nil {
		return err
	}
	fmt.Println("WE DOIN IT NOW")

	req, err := http.NewRequest(http.MethodPost, u.String(), data)
	if err != nil {
		return err
	}
	access, err := getAccessToken()
	if err != nil {
		return err
	}
	fmt.Println("GOT THAT TOKEN")
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

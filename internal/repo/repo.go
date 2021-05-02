package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

type repoService interface {
	DownloadPlugin(plugin *api.Plugin) ([]byte, error)
	UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error
}

type repoServiceImpl struct {
	Host string
	Port string
}

func NewRepoService() repoService {
	return &repoServiceImpl{
		Host: os.Getenv("REPO_HOST"),
		Port: os.Getenv("REPO_PORT"),
	}
}

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", gate.SimpleAuth(authUpload(pluginsHandler)))
	mux.Handle("/repo/thumbnails", gate.SimpleAuth(authUpload(thumbnailsHandler)))

	return mux
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
	access, err := gate.GetAccessToken()
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
	access, err := gate.GetAccessToken()
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

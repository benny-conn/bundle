package gate

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

type gateService interface {
	DownloadPlugin(plugin *api.Plugin) ([]byte, error)
	UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error
	PaginatePlugins(page int, count int) ([]*api.Plugin, error)
	GetPlugin(plugin *api.Plugin) (*api.Plugin, error)
	InsertPlugin(plugin *api.Plugin) error
	UpdatePlugin(updatedPlugin *api.Plugin) error
	GetReadme(plugin *api.Plugin) (*api.Readme, error)
	InsertReadme(user *api.User, readme *api.Readme) error
	UpdateReadme(user *api.User, readme *api.Readme) error
	UpdateUser(updatedUser *api.User) error
	GetUser(user *api.User) (*api.User, error)
	InsertUser(user *api.User) error
}
type gateServiceImpl struct {
	Host string
	Port string
}

func NewGateService(host string, port string) gateService {
	if host == "" {
		host = os.Getenv("GATE_HOST")
	}
	if port == "" {
		port = os.Getenv("GATE_PORT")
	}
	return &gateServiceImpl{
		Host: host,
		Port: port,
	}
}

func (g *gateServiceImpl) DownloadPlugin(plugin *api.Plugin) ([]byte, error) {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/repo/plugins", g.Host, g.Port))
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

func (g *gateServiceImpl) UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/repo/plugins", g.Host, g.Port))
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

func (g *gateServiceImpl) PaginatePlugins(page int, count int) ([]*api.Plugin, error) {
	addr := fmt.Sprintf("http://%v:%v/api/plugins", g.Host, g.Port)
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

	// FIGURE OUT WHAT THE HECK TO DO HERE
	result := &api.PaginatePluginsResponse{}
	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}
	return result.Plugins, nil

}

func (g *gateServiceImpl) GetPlugin(plugin *api.Plugin) (*api.Plugin, error) {

	addr := fmt.Sprintf("http://%v:%v/api/plugins", g.Host, g.Port)
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

func (g *gateServiceImpl) InsertPlugin(plugin *api.Plugin) error {

	addr := fmt.Sprintf("http://%v:%v/api/plugins", g.Host, g.Port)
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

func (g *gateServiceImpl) UpdatePlugin(updatedPlugin *api.Plugin) error {

	addr := fmt.Sprintf("http://%v:%v/api/plugins", g.Host, g.Port)
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

func (g *gateServiceImpl) InsertReadme(user *api.User, readme *api.Readme) error {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/readmes", g.Host, g.Port))
	if err != nil {
		return err
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	readmeJSON, err := json.Marshal(readme)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte(readmeJSON))

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Add("User", string(userJSON))
	req.Header.Add("Resource", string(readmeJSON))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (g *gateServiceImpl) GetReadme(plugin *api.Plugin) (*api.Readme, error) {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/readmes", g.Host, g.Port))
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

	rdme := &api.Readme{}
	err = json.Unmarshal(bs, &rdme)
	if err != nil {
		return nil, err
	}
	return rdme, nil
}

func (g *gateServiceImpl) UpdateReadme(user *api.User, readme *api.Readme) error {

	u, err := url.Parse(fmt.Sprintf("http://%v:%v/api/readmes", g.Host, g.Port))
	if err != nil {
		return err
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	readmeJSON, err := json.Marshal(readme)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte(readmeJSON))

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Add("User", string(userJSON))
	req.Header.Add("Resource", string(readmeJSON))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (g *gateServiceImpl) UpdateUser(updatedUser *api.User) error {

	addr := fmt.Sprintf("http://%v:%v/api/users", g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	updatedBs, err := json.Marshal(updatedUser)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(updatedBs)

	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	access, err := getAccessToken()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+access)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (g *gateServiceImpl) GetUser(user *api.User) (*api.User, error) {

	addr := fmt.Sprintf("http://%v:%v/api/users", g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("id", user.Id)
	q.Set("username", user.Username)
	q.Set("email", user.Email)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	access, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.User{}

	err = json.Unmarshal(bs, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *gateServiceImpl) InsertUser(user *api.User) error {

	addr := fmt.Sprintf("http://%v:%v/api/users", g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	bs, err := json.Marshal(user)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	buf.Write(bs)

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := getAccessToken()
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+access)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

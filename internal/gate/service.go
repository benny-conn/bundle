package gate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
)

type gateService interface {
	DownloadPlugin(plugin *api.Plugin) ([]byte, error)
	UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error
	UploadThumbnail(user *api.User, plugin *api.Plugin, data io.Reader) error
	PaginatePlugins(req *api.PaginatePluginsRequest) ([]*api.Plugin, error)
	GetPlugin(plugin *api.Plugin) (*api.Plugin, error)
	InsertPlugin(plugin *api.Plugin) error
	UpdatePlugin(updatedPlugin *api.Plugin) error
	GetReadme(plugin *api.Plugin) (*api.Readme, error)
	InsertReadme(user *api.User, readme *api.Readme) error
	UpdateReadme(user *api.User, readme *api.Readme) error
	UpdateUser(updatedUser *api.User) error
	GetUser(user *api.User) (*api.User, error)
	InsertUser(user *api.User) error
	InsertSession(ses *api.Session) (*api.SessionInsertResponse, error)
	GetSession(ses *api.Session) (*api.Session, error)
	DeleteSession(ses *api.Session) error
	InsertBundle(bu *api.Bundle) error
	GetBundle(bu *api.Bundle) (*api.Bundle, error)
	DeleteBundle(bu *api.Bundle) error
	UpdateBundle(bu *api.Bundle) error
	GetChangelog(ch *api.Changelog) (*api.Changelog, error)
	InsertChangelog(ch *api.Changelog) error
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

	scheme := "https://"
	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/repo/plugins", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", plugin.Name)
	q.Add("version", plugin.Version)
	u.RawQuery = q.Encode()

	client := internal.NewBasicClient()

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}
	return bs.Bytes(), nil
}

func (g *gateServiceImpl) UploadPlugin(user *api.User, plugin *api.Plugin, data io.Reader) error {
	scheme := "https://"
	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/repo/plugins", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if user != nil && plugin != nil {
		if user.Username == "" || user.Password == "" || plugin.Name == "" || plugin.Version == "" {
			return errors.New("missing required fields")
		}
	} else {
		return errors.New("specify a user and plugin")
	}

	writer.WriteField("username", user.Username)
	writer.WriteField("author", user.Username)
	writer.WriteField("password", user.Password)
	writer.WriteField("name", plugin.Name)
	writer.WriteField("version", plugin.Version)
	writer.WriteField("description", plugin.Description)
	writer.WriteField("category", fmt.Sprint(plugin.Category))

	part, err := writer.CreateFormFile("plugin", plugin.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, data)
	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	client := internal.NewBasicClient()
	resp, err := client.Post(u.String(), writer.FormDataContentType(), body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) UploadThumbnail(user *api.User, plugin *api.Plugin, data io.Reader) error {
	scheme := "https://"
	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/repo/thumbnails", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if user == nil && plugin == nil {
		return errors.New("specify a user or plugin")
	}
	if plugin != nil {
		if plugin.Id == "" {
			return errors.New("specify a plugin id")
		}
		writer.WriteField("plugin", plugin.Id)
	}
	if user != nil {
		if user.Id == "" {
			return errors.New("specify a user id")
		}
		writer.WriteField("user", user.Id)
	}

	part, err := writer.CreateFormFile("thumbnail", "THUMBNAIL.webp")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, data)
	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodPost, u.String(), body)

	if err != nil {
		return err
	}

	accessToken, err := newAuthToken("thumbnails")
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) PaginatePlugins(req *api.PaginatePluginsRequest) ([]*api.Plugin, error) {
	scheme := "https://"
	addr := fmt.Sprintf("%s%s:%s/api/plugins", scheme, g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", fmt.Sprint(req.Page))
	q.Set("count", fmt.Sprint(req.Count))
	q.Set("search", req.Search)
	q.Set("category", fmt.Sprint(req.Category))
	q.Set("sort", fmt.Sprint(req.Sort))
	u.RawQuery = q.Encode()

	client := internal.NewBasicClient()
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.PaginatePluginsResponse{}
	err = json.Unmarshal(bs.Bytes(), result)
	if err != nil {
		return nil, err
	}
	return result.Plugins, nil

}

func (g *gateServiceImpl) GetPlugin(plugin *api.Plugin) (*api.Plugin, error) {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/plugins", scheme, g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("name", plugin.Name)
	q.Set("id", plugin.Id)
	u.RawQuery = q.Encode()

	client := internal.NewBasicClient()
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.Plugin{}
	err = json.Unmarshal(bs.Bytes(), result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *gateServiceImpl) InsertPlugin(plugin *api.Plugin) error {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/plugins", scheme, g.Host, g.Port)
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

	client := internal.NewBasicClient()
	resp, err := client.Post(u.String(), "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) UpdatePlugin(updatedPlugin *api.Plugin) error {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/plugins", scheme, g.Host, g.Port)
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

	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) InsertReadme(user *api.User, readme *api.Readme) error {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/readmes", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}

	client := internal.NewBasicClient()

	values := url.Values{}
	values.Set("username", user.Username)
	values.Set("password", user.Password)
	values.Set("plugin_id", readme.Plugin.Id)
	values.Set("plugin_name", readme.Plugin.Name)
	values.Set("text", readme.Text)

	resp, err := client.PostForm(u.String(), values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) GetReadme(plugin *api.Plugin) (*api.Readme, error) {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/readmes", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("name", plugin.Name)
	q.Add("id", plugin.Id)
	u.RawQuery = q.Encode()
	client := internal.NewBasicClient()
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	rdme := &api.Readme{}
	err = json.Unmarshal(bs.Bytes(), &rdme)
	if err != nil {
		return nil, err
	}
	return rdme, nil
}

func (g *gateServiceImpl) UpdateReadme(user *api.User, readme *api.Readme) error {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/readmes", scheme, g.Host, g.Port))
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
	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Add("User", string(userJSON))
	req.Header.Add("Resource", string(readmeJSON))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) UpdateUser(updatedUser *api.User) error {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/users", scheme, g.Host, g.Port)
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
	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}

	access, err := newAuthToken("users")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) GetUser(user *api.User) (*api.User, error) {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/users", scheme, g.Host, g.Port)
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("id", user.Id)
	q.Set("username", user.Username)
	q.Set("email", user.Email)
	u.RawQuery = q.Encode()

	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	access, err := newAuthToken("users")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.User{}

	err = json.Unmarshal(bs.Bytes(), result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *gateServiceImpl) InsertUser(user *api.User) error {
	scheme := "https://"

	addr := fmt.Sprintf("%s%s:%s/api/users", scheme, g.Host, g.Port)
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

	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := newAuthToken("users")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) DeleteSession(ses *api.Session) error {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/sessions", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}
	asJSON, err := json.Marshal(ses)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(asJSON))
	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodDelete, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := newAuthToken("sessions")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}
	return nil
}

func (g *gateServiceImpl) GetSession(ses *api.Session) (*api.Session, error) {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/sessions", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("id", ses.Id)
	q.Add("userId", ses.UserId)
	u.RawQuery = q.Encode()
	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	access, err := newAuthToken("sessions")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.Session{}
	err = json.Unmarshal(bs.Bytes(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *gateServiceImpl) InsertSession(ses *api.Session) (*api.SessionInsertResponse, error) {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/sessions", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}

	asJSON, err := json.Marshal(ses)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte(asJSON))

	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return nil, err
	}
	access, err := newAuthToken("sessions")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}

	if _, err = io.Copy(bs, resp.Body); err != nil {
		return nil, err
	}

	result := &api.SessionInsertResponse{}

	if err = json.Unmarshal(bs.Bytes(), result); err != nil {
		return nil, err
	}

	return result, nil
}

func (g *gateServiceImpl) DeleteBundle(bu *api.Bundle) error {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/bundles", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}
	asJSON, err := json.Marshal(bu)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(asJSON))
	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodDelete, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := newAuthToken("bundles")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}

	return nil
}

func (g *gateServiceImpl) UpdateBundle(bu *api.Bundle) error {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/bundles", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}
	asJSON, err := json.Marshal(bu)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(asJSON))
	client := internal.NewBasicClient()
	req, err := http.NewRequest(http.MethodPatch, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := newAuthToken("bundles")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}
	return nil
}

func (g *gateServiceImpl) GetBundle(bu *api.Bundle) (*api.Bundle, error) {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/bundles", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("id", bu.Id)
	q.Add("userId", bu.UserId)
	u.RawQuery = q.Encode()
	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	access, err := newAuthToken("bundles")
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+access)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.Bundle{}
	err = json.Unmarshal(bs.Bytes(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *gateServiceImpl) InsertBundle(bu *api.Bundle) error {

	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/bundles", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}

	asJSON, err := json.Marshal(bu)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(asJSON))

	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	access, err := newAuthToken("bundles")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+access)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}
	return nil
}

func (g *gateServiceImpl) InsertChangelog(ch *api.Changelog) error {

	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/changelogs", scheme, g.Host, g.Port))
	if err != nil {
		return err
	}

	asJSON, err := json.Marshal(ch)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(asJSON))

	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return err
		}
		return errors.New(buf.String())
	}
	return nil
}

func (g *gateServiceImpl) GetChangelog(ch *api.Changelog) (*api.Changelog, error) {
	scheme := "https://"

	u, err := url.Parse(fmt.Sprintf("%s%s:%s/api/changelogs", scheme, g.Host, g.Port))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("id", ch.Id)
	q.Add("pluginId", ch.PluginId)
	q.Add("version", ch.Version)
	u.RawQuery = q.Encode()
	client := internal.NewBasicClient()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if internal.IsRespError(resp) {
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(buf.String())
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, resp.Body)
	if err != nil {
		return nil, err
	}

	result := &api.Changelog{}
	err = json.Unmarshal(bs.Bytes(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

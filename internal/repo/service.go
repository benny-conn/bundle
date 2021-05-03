package repo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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
	q.Add("id", plugin.Id)
	q.Add("author", plugin.Author.Id)
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.SetBoundary("XXX")

	if plugin.Id == "" || plugin.Version == "" || plugin.Author.Id == "" {
		return errors.New("missing required fields")
	}

	writer.WriteField("id", plugin.Id)
	writer.WriteField("version", plugin.Version)
	writer.WriteField("author", plugin.Author.Id)

	part, err := writer.CreateFormFile("plugin", plugin.Id)
	if err != nil {
		return err
	}
	bs, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	_, err = part.Write(bs)
	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	access, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+access)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=XXX")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))

	defer resp.Body.Close()
	return nil
}

package pkg

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/bennycio/bundle/api"
)

func PaginatePlugins(page int) ([]api.Plugin, error) {

	u, err := url.Parse(ApiServerHost)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", string(page))
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

func GetPlugin(name string) (*api.Plugin, error) {

	u, err := url.Parse(ApiServerHost + "/plugins")
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

package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

func GetPluginData(opts *api.GetPluginDataRequest) (*api.Plugin, error) {
	u, err := url.Parse(ApiServerHost + "/bundles")
	if err != nil {
		return nil, err
	}

	withPlugin := strconv.FormatBool(opts.WithPlugin)
	withReadme := strconv.FormatBool(opts.WithReadme)
	withThumbnail := strconv.FormatBool(opts.WithThumbnail)

	q := u.Query()
	q.Set("name", opts.Name)
	q.Set("version", opts.Version)
	q.Set("withPlugin", withPlugin)
	q.Set("withReadme", withReadme)
	q.Set("withThumbnail", withThumbnail)
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

func InsertPluginData(req *api.InsertPluginDataRequest) error {
	u, err := url.Parse(ApiServerHost + "/bundles")
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	switch req.DataType.String() {
	case api.InsertPluginDataRequest_PLUGIN.String():
		_, err = buf.Write(req.Plugin.PluginData)
		if err != nil {
			return err
		}
	case api.InsertPluginDataRequest_README.String():
		_, err = buf.Write(req.Plugin.Readme)
		if err != nil {
			return err
		}
	case api.InsertPluginDataRequest_THUMBNAIL.String():
		_, err = buf.Write(req.Plugin.Thumbnail)
		if err != nil {
			return err
		}
	default:
		_, err = buf.Write(req.Plugin.PluginData)
		if err != nil {
			return err
		}
	}
	req.Plugin.PluginData = nil
	req.Plugin.Readme = nil
	req.Plugin.Thumbnail = nil

	request, err := http.NewRequest(http.MethodPost, u.String(), buf)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(req.Plugin)
	if err != nil {
		return err
	}
	request.Header.Add("Data-Type", req.DataType.String())
	request.Header.Add("Resource", string(bs))

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil

}

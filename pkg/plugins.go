package pkg

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bennycio/bundle/internal"
)

func GetPlugin(pluginName string) (*internal.Plugin, error) {
	resp, err := http.Get("http://localhost:8080/plugins?plugin=" + pluginName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &internal.Plugin{}

	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

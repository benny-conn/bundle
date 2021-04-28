package wrapper

import (
	"io"

	"github.com/bennycio/bundle/api"
)

func DownloadReadme(pluginName string) (string, error) {}

func DownloadPlugin(pluginName string, version string) ([]byte, error) {}

func UploadPlugin(user *api.User, pluginName string, version string, data io.Reader) error {}

func UploadReadme(user *api.User, pluginName string, data io.Reader) error {}

func DownloadReadmeApi(pluginName string) (string, error) {}

func DownloadPluginApi(pluginName string, version string) ([]byte, error) {}

func UploadPluginApi(user *api.User, pluginName string, version string, data io.Reader) error {}

func UploadReadmeApi(user *api.User, pluginName string, data io.Reader) error {}

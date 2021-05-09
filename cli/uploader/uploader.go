package uploader

import (
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

type uploader struct {
	location string
	user     *api.User
	plugin   *api.Plugin
	opts     options
}

type options struct {
	isReadme bool
}

func New(user *api.User, location, name, version string) *uploader {
	return &uploader{
		user:     user,
		location: location,
		plugin: &api.Plugin{
			Name:    name,
			Version: version,
		},
	}
}

func (u *uploader) WithReadme(isReadme bool) *uploader {
	u.opts.isReadme = isReadme
	return u
}

func (u *uploader) Upload() error {
	gservice := gate.NewGateService("localhost", "8020")
	if u.opts.isReadme {
		file, err := os.ReadFile(u.location)
		if err != nil {
			return err
		}
		readme := &api.Readme{
			Plugin: u.plugin,
			Text:   string(file),
		}
		err = gservice.InsertReadme(u.user, readme)
		if err != nil {
			return err
		}
	} else {
		file, err := os.Open(u.location)
		if err != nil {
			return err
		}
		defer file.Close()
		err = gservice.UploadPlugin(u.user, u.plugin, file)
		if err != nil {
			return err
		}
	}
	return nil
}

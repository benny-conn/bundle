package uploader

import (
	"bytes"
	"io"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

type uploader struct {
	file   *os.File
	user   *api.User
	plugin *api.Plugin
	opts   options
}

type options struct {
	isReadme bool
}

func New(user *api.User, file *os.File, plugin *api.Plugin) *uploader {
	return &uploader{
		user:   user,
		file:   file,
		plugin: plugin,
	}
}

func (u *uploader) WithReadme(isReadme bool) *uploader {
	u.opts.isReadme = isReadme
	return u
}

func (u *uploader) Upload() error {
	gservice := gate.NewGateService("localhost", "8020")
	if u.opts.isReadme {
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, u.file)
		if err != nil {
			return err
		}
		readme := &api.Readme{
			Plugin: u.plugin,
			Text:   buf.String(),
		}
		err = gservice.InsertReadme(u.user, readme)
		if err != nil {
			return err
		}
	} else {
		err := gservice.UploadPlugin(u.user, u.plugin, u.file)
		if err != nil {
			return err
		}
	}
	return nil
}

package uploader

import (
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/schollz/progressbar/v3"
)

type Uploader struct {
	PluginFile *os.File
	User       *api.User
	Plugin     *api.Plugin
	Readme     *api.Readme
	Changelog  *api.Changelog
}

func (u *Uploader) Upload() error {
	gservice := gate.NewGateService("localhost", "8020")
	if u.Plugin != nil && u.PluginFile != nil {
		fi, err := u.PluginFile.Stat()
		if err != nil {
			return err
		}
		pb := progressbar.DefaultBytes(fi.Size(), "Uploading Plugin...")

		rdr := progressbar.NewReader(u.PluginFile, pb)

		err = gservice.UploadPlugin(u.User, u.Plugin, &rdr)
		if err != nil {
			return err
		}
	}
	if u.Readme != nil {

		pb := progressbar.Default(1, "Uploading Readme...")
		if err := gservice.InsertReadme(u.User, u.Readme); err != nil {
			return err
		}
		pb.Add(1)
	}
	if u.Changelog != nil {
		pb := progressbar.Default(1, "Uploading Changelog...")
		if err := gservice.InsertChangelog(u.User, u.Changelog); err != nil {
			return err
		}
		pb.Add(1)
	}
	return nil
}

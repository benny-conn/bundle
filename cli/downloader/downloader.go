package downloader

import (
	"bytes"
	"io"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
)

type downloader struct {
	Plugin  *api.Plugin
	options options
}

type options struct {
	latest   bool
	location string
}

func New(name, version string) *downloader {
	return &downloader{
		Plugin: &api.Plugin{
			Name:    name,
			Version: version,
		},
	}
}

func (d *downloader) WithLatest(latest bool) *downloader {
	d.options.latest = latest
	return d
}

func (d *downloader) WithLocation(loc string) *downloader {
	d.options.location = loc
	return d
}

func (u *downloader) Download() ([]byte, error) {

	gservice := gate.NewGateService("localhost", "8020")

	if u.options.latest {
		plugin, err := gservice.GetPlugin(u.Plugin)
		if err != nil {

			return nil, err
		}
		u.Plugin.Version = plugin.Version
	}

	pl, err := gservice.DownloadPlugin(u.Plugin)
	if err != nil {

		return nil, err
	}

	return pl, nil
}

func (d *downloader) Install(bs []byte) error {

	file, err := os.OpenFile(d.options.location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {

		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bs)

	_, err = io.Copy(file, buf)

	if err != nil {
		return err
	}

	return nil

}

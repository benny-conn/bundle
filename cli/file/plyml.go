package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jlaffaye/ftp"
	"gopkg.in/yaml.v2"
)

type PluginYml struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description,omitempty"`
	Category    int32    `yaml:"category,omitempty"`
	Conflicts   []string `yaml:"conflicts,omitempty"`
}

func ParsePluginYml(rd io.ReaderAt, size int64) (PluginYml, error) {

	reader, err := zip.NewReader(rd, size)

	if err != nil {
		return PluginYml{}, err
	}

	result := PluginYml{}

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "plugin.yml") {
			rc, err := file.Open()
			if err != nil {
				return PluginYml{}, err
			}
			buf := bytes.Buffer{}
			buf.ReadFrom(rc)
			err = yaml.Unmarshal(buf.Bytes(), &result)
			if err != nil {
				return PluginYml{}, err
			}
		}
	}

	return result, nil
}

func GetPluginYml(pluginName string, conn *ftp.ServerConn) (PluginYml, error) {

	if conn == nil {
		plfile, err := os.Open(fmt.Sprintf("plugins/%s.jar", pluginName))
		if err == nil {
			defer plfile.Close()

			info, err := plfile.Stat()
			if err != nil {
				return PluginYml{}, err
			}
			p, err := ParsePluginYml(plfile, info.Size())
			if err != nil {
				return PluginYml{}, err
			}
			return p, nil
		} else {
			return PluginYml{}, err
		}
	} else {
		fp := fmt.Sprintf("plugins/%s.jar", pluginName)
		resp, err := conn.Retr(fp)
		if err == nil {
			defer resp.Close()
			buf := &bytes.Buffer{}
			size, err := io.Copy(buf, resp)
			if err != nil {
				return PluginYml{}, err
			}
			rderAt := bytes.NewReader(buf.Bytes())
			p, err := ParsePluginYml(rderAt, size)
			if err != nil {
				return PluginYml{}, err
			}
			return p, nil
		} else {

			return PluginYml{}, err
		}
	}

}

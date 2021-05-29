package file

import (
	"archive/zip"
	"bytes"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type PlYml struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description,omitempty"`
	Category    int32    `yaml:"category,omitempty"`
	Conflicts   []string `yaml:"conflicts,omitempty"`
}

func ParsePluginYml(file *os.File) (PlYml, error) {

	info, err := file.Stat()
	if err != nil {
		return PlYml{}, err
	}
	reader, err := zip.NewReader(file, info.Size())

	if err != nil {
		return PlYml{}, err
	}

	result := PlYml{}

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "plugin.yml") {
			rc, err := file.Open()
			if err != nil {
				return PlYml{}, err
			}
			buf := bytes.Buffer{}
			buf.ReadFrom(rc)
			err = yaml.Unmarshal(buf.Bytes(), &result)
			if err != nil {
				return PlYml{}, err
			}
		}
	}

	return result, nil

}

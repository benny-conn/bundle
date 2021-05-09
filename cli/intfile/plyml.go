package intfile

import (
	"archive/zip"
	"bytes"
	"strings"

	"gopkg.in/yaml.v2"
)

type plYml struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`
}

func Parse(path string) (plYml, error) {
	reader, err := zip.OpenReader(path)

	if err != nil {
		return plYml{}, err
	}
	defer reader.Close()

	result := plYml{}

	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "plugin.yml") {
			rc, err := file.Open()
			if err != nil {
				return plYml{}, err
			}
			buf := bytes.Buffer{}
			buf.ReadFrom(rc)
			err = yaml.Unmarshal(buf.Bytes(), &result)
			if err != nil {
				return plYml{}, err
			}
		}
	}

	return result, nil

}

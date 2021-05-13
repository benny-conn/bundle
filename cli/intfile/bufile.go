package intfile

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"gopkg.in/yaml.v2"
)

const (
	BuFileName = "bundle.yml"
)

//go:embed bundle.yml
var buFile string

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

func Initialize(path string) error {
	file, err := os.Create(filepath.Join(path, BuFileName))
	if err != nil {
		return err
	}
	_, err = file.WriteString(buFile)
	if err != nil {
		return err
	}
	return nil
}

func IsBundleInitialized(path string) bool {
	path = findBundle(path)
	_, err := os.Stat(path)
	return err == nil
}

func getBundleFileBytes(path string) ([]byte, error) {

	if !IsBundleInitialized(path) {
		return nil, errors.New("bundle file does not exist at current directory")
	}
	path = findBundle(path)

	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bs := &bytes.Buffer{}
	_, err = io.Copy(bs, fi)
	if err != nil {
		panic(err)
	}
	return bs.Bytes(), nil
}

func GetBundle(path string) (*BundleFile, error) {

	fileBytes, err := getBundleFileBytes(path)

	if err != nil {
		return nil, err
	}
	result := &BundleFile{}

	err = yaml.Unmarshal(fileBytes, result)

	if err != nil {
		panic(err)
	}

	return result, nil
}

func WritePluginsToBundle(plugins map[string]string, path string) error {
	bundle, err := GetBundle(path)
	if err != nil {
		return err
	}
	bundle.Plugins = plugins
	path = findBundle(path)

	currentBundleFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	currentBundleFile.Truncate(0)
	newFileBytes, err := yaml.Marshal(bundle)
	if err != nil {
		return err
	}
	_, err = currentBundleFile.Write(newFileBytes)
	if err != nil {
		return err
	}
	return nil
}

func findBundle(path string) string {
	result := path
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	if path == "" || !strings.HasSuffix(path, BuFileName) {
		result = filepath.Join(wd, BuFileName)
	}
	return result
}

package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/jlaffaye/ftp"
	"gopkg.in/yaml.v2"
)

const (
	BuFileName = "bundle.yml"
)

//go:embed bundle.yml
var BuFile string

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

func Initialize(path string) error {
	file, err := os.Create(filepath.Join(path, BuFileName))
	if err != nil {
		return err
	}
	_, err = file.WriteString(BuFile)
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

func GetBundleFtp(conn *ftp.ServerConn) (BundleFile, error) {
	resp, err := conn.Retr("bundle.yml")
	if err != nil {
		return BundleFile{}, err
	}
	defer resp.Close()
	result := &BundleFile{}

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp)
	if err != nil {
		return BundleFile{}, err
	}

	err = yaml.Unmarshal(buf.Bytes(), result)
	if err != nil {
		return BundleFile{}, err
	}
	return *result, nil
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

func WritePluginsToBundleFtp(conn *ftp.ServerConn, plugins map[string]string) error {
	bundle, err := GetBundleFtp(conn)
	if err != nil {
		return err
	}
	bundle.Plugins = plugins

	conn.Delete(BuFileName)
	newFileBytes, err := yaml.Marshal(bundle)
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		_, err := pw.Write(newFileBytes)
		if err != nil {
			fmt.Printf("error occurred: %s\n", err.Error())
			return
		}
	}()

	if err := conn.Stor(BuFileName, pr); err != nil {
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

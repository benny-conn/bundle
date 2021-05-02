package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bennycio/bundle/api"
	"gopkg.in/yaml.v2"
)

func isBundleInitialized() bool {
	fn := BundleFileName
	_, err := os.Stat(fn)
	return err == nil
}

func getBundleFileBytes() ([]byte, error) {

	if !isBundleInitialized() {
		return nil, errors.New("bundle file does not exist at current directory")
	}
	fn := BundleFileName

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}

func isPluginDirectory() bool {
	if _, err := os.Stat("plugins/"); os.IsNotExist(err) {
		return false
	}
	return true
}

func credentialsPrompt() *api.User {

	fmt.Println("Enter your username: ")
	var username string
	fmt.Scanln(&username)
	fmt.Println("Enter your password: ")
	var password string
	fmt.Scanln(&password)

	user := &api.User{}
	user.Password = password
	user.Username = username

	return user
}

func getBundleFilePlugins() (map[string]string, error) {

	result, err := getBundle()
	if err != nil {
		return nil, err
	}
	return result.Plugins, nil
}

func getBundle() (*BundleFile, error) {

	fileBytes, err := getBundleFileBytes()

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

func writePluginsToBundle(plugins map[string]string) error {
	bundle, err := getBundle()
	if err != nil {
		return err
	}
	bundle.Plugins = plugins

	currentBundleFile, err := os.OpenFile(BundleFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

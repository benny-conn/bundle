package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	bundle "github.com/bennycio/bundle/internal"
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

func credentialsPrompt() *bundle.User {

	fmt.Println("Enter your username or email: ")
	var userOrEmail string
	fmt.Scanln(&userOrEmail)
	fmt.Println("Enter your password: ")
	var password string
	fmt.Scanln(&password)

	isEmail := emailRegex.MatchString(userOrEmail)

	user := &bundle.User{}
	user.Password = password
	if isEmail {
		user.Email = userOrEmail
	} else {
		user.Username = userOrEmail
	}

	return user
}

func getBundleFilePlugins() (map[string]string, error) {

	fileBytes, err := getBundleFileBytes()

	if err != nil {
		return nil, err
	}
	result := &BundleFile{}

	err = yaml.Unmarshal(fileBytes, result)

	if err != nil {
		return nil, err
	}

	return result.Plugins, nil
}

func getBundleFile() (*BundleFile, error) {

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

func updateBundleVersion(pluginName string, version string) error {
	file, err := getBundleFile()
	if err != nil {
		return err
	}
	file.Plugins[pluginName] = version

	// TODO adjust the file itself

	return nil
}

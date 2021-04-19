package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	bundle "github.com/bennycio/bundle/internal"
	"gopkg.in/yaml.v2"
)

func getPlugin(pluginName string) (*bundle.Plugin, error) {
	resp, err := http.Get("http://localhost:8080/plugins?plugin=" + pluginName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &bundle.Plugin{}

	err = json.Unmarshal(bs, result)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func isBundleInitialized() bool {
	fn := BundleFileName
	_, err := os.Stat(fn)
	return err == nil
}

func getBundleFile() ([]byte, error) {

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

func getBundledPlugins() map[string]string {
	fileBytes, err := getBundleFile()

	if err != nil {
		panic(err)
	}
	result := BundleFile{}

	err = yaml.Unmarshal(fileBytes, &result)

	if err != nil {
		panic(err)
	}

	return result.Plugins
}

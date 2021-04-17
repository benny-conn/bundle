package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

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

func credentialsPrompt() *User {

	fmt.Println("Enter your username or email: ")
	var userOrEmail string
	fmt.Scanln(&userOrEmail)
	fmt.Println("Enter your password: ")
	var password string
	fmt.Scanln(&password)

	isEmail := emailRegex.MatchString(userOrEmail)

	user := &User{}
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

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bennycio/bundle/api"
)

func isPluginDirectory(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "plugins/")); os.IsNotExist(err) {
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

	user := &api.User{
		Username: username,
		Password: password,
	}

	return user
}

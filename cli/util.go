package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bennycio/bundle/api"
	"golang.org/x/crypto/ssh/terminal"
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

	fmt.Print("Enter Your Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err.Error())
	}
	password := string(bytePassword)

	user := &api.User{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	return user
}

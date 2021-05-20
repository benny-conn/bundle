package cli

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/term"
	"golang.org/x/crypto/ssh/terminal"
)

func isPluginDirectory(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "plugins/")); os.IsNotExist(err) {
		return false
	}
	return true
}

func credentialsPrompt() *api.User {

	term.Println("Enter your username: ")
	rd := bufio.NewReader(os.Stdin)
	username, _ := rd.ReadString(byte('\n'))
	username = strings.TrimSpace(strings.Trim(username, "\n"))

	term.Println("Enter Your Password: ")
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

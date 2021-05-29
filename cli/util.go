package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/term"
	"github.com/c-bata/go-prompt"
	goterm "golang.org/x/term"
)

func isPluginDirectory(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "plugins/")); os.IsNotExist(err) {
		return false
	}
	return true
}

func credentialsPrompt() *api.User {

	term.Println("Enter your username: ")

	username := prompt.Input(">> ", nilCompleter)

	term.Println("Enter Your Password: ")
	bytePassword, err := goterm.ReadPassword(syscall.Stdin)
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

func nilCompleter(d prompt.Document) []prompt.Suggest {
	return nil
}

func versionGreaterThan(version, than string) bool {

	isNumber := func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return 0
	}

	versionNoChars := strings.Map(isNumber, version)

	thanNoChars := strings.Map(isNumber, than)

	splitInts := make([]int, len(versionNoChars))
	thanInts := make([]int, len(thanNoChars))

	for i, v := range versionNoChars {
		in := int(v - '0')
		splitInts[i] = in
	}

	for i, v := range thanNoChars {
		in := int(v - '0')
		thanInts[i] = in
	}

	for i, v := range thanInts {
		if len(versionNoChars) < i+1 {
			break
		}
		if splitInts[i] == v {
			continue
		}
		return splitInts[i] > v
	}

	return false

}

func completerWithOptions(ss ...string) func(prompt.Document) []prompt.Suggest {
	suggests := []prompt.Suggest{}

	for _, v := range ss {
		suggests = append(suggests, prompt.Suggest{Text: v})
	}
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggests, d.GetWordBeforeCursor(), true)
	}
}

func yesOrNoCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "y"},
		{Text: "yes"},
		{Text: "n"},
		{Text: "no"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

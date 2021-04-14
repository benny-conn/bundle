package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	_ "embed"
)

const (
	FILE_NAME          string = "bundle.yml"
	MAKE_FILE_NAME     string = "bundle-make.yml"
	REQUIRED_FILE_TYPE string = "application/zip"
)

type BundleFile struct {
	Plugins map[string]string `yaml:"Plugins,omitempty"`
}

type BundleMakeFile struct {
	Name         string   `yaml:"Name,omitempty"`
	Version      string   `yaml:"Version,omitempty"`
	JarPath      string   `yaml:"JarPath,omitempty"`
	VersionNotes []string `yaml:"VersionNotes,omitempty"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//go:embed bundle.yml
var BundleYml string

//go:embed bundle-make.yml
var BundleMakeYml string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO
}

func isBundleInitialized(make bool) bool {
	var fn string
	if make {
		fn = MAKE_FILE_NAME
	} else {
		fn = FILE_NAME
	}
	_, err := os.Stat(fn)
	return err == nil
}

func getBundleFile(make bool) ([]byte, error) {

	var fn string

	if make {
		fn = MAKE_FILE_NAME
	} else {
		fn = FILE_NAME
	}
	if !isBundleInitialized(make) {
		return nil, errors.New("bundle file does not exist at current directory")
	}

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
	fileBytes, err := getBundleFile(false)

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

package cli

import (
	"regexp"

	"github.com/spf13/cobra"

	_ "embed"

	bundle "github.com/bennycio/bundle/internal"
)

const (
	BundleFileName   = bundle.BundleFileName
	RequiredFileType = bundle.RequiredFileType
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//go:embed bundle.yml
var BundleYml string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Base command for the Bundle CLI",
}

var Force bool

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Force the command to run regardless of constraints")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO
}

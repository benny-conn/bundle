package cli

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Base command for the Bundle CLI",
}

var force bool

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Specify a path to bundle file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO
}

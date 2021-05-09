package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	_ "embed"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Base command for the Bundle CLI",
}

var force bool

var buFilePath string

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not find working directory")
	}
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force the command to run regardless of constraints")
	rootCmd.PersistentFlags().StringVarP(&buFilePath, "bundle", "b", wd, "Specify a path to bundle file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO
}

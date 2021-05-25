package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Base command for the Bundle CLI",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yml")           // REQUIRED if the config file does not have the extension in the name  // path to look for the config file in
	viper.AddConfigPath("$HOME/.bundle") // call multiple times to add many search paths            // optionally look for config in the working directory
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %s", err))
	}
}

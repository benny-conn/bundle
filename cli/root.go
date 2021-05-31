package cli

import (
	"fmt"
	"os"

	"github.com/bennycio/bundle/cli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configPath string

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
	confDir, err := os.UserConfigDir()
	if err != nil {
		confDir, err = os.UserHomeDir()
		if err != nil {
			logger.ErrLog.Fatal(err.Error())
		}
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(fmt.Sprintf("%s/.bundle", confDir))
	viper.SetDefault("FTP", map[string]map[string]string{})
	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fmt.Sprintf("%s/.bundle", confDir), os.ModePerm)
			if err != nil {
				logger.ErrLog.Fatal(err.Error())
			}
			err = viper.WriteConfigAs(fmt.Sprintf("%s/.bundle/config.yml", confDir))
			if err != nil {
				logger.ErrLog.Fatal(err.Error())
			}
		}
	}
	err = viper.ReadInConfig()
	if err != nil {
		logger.ErrLog.Fatal(err.Error())
	}

	configPath = fmt.Sprintf("%s/.bundle/config.yml", confDir)

}

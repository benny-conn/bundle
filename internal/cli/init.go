package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup structure for running a server with access to the official Bundle Repository",
	Run: func(cmd *cobra.Command, args []string) {

		if isBundleInitialized(false) {
			log.Fatal("There already exists a bundle.yml at this location")
		}

		if !isPluginDirectory() {
			log.Fatal("There is no plugin directory in your current directory")
		}

		file, err := os.Create(FILE_NAME)
		if err != nil {
			panic(err)
		}
		file.WriteString(BundleYml)

		wd, _ := os.Getwd()

		fmt.Println("Created file at path " + wd + "/" + FILE_NAME)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

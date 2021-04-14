package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// makeCmd represents the create command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Setup structure for uploading plugins to the official Bundle Repository",
	Run: func(cmd *cobra.Command, args []string) {

		if isBundleInitialized(true) {
			log.Fatal("There already exists a bundle-make.yml at this location")
		}

		file, err := os.Create(BundleMakeFileName)
		if err != nil {
			panic(err)
		}
		file.WriteString(BundleMakeYml)

		wd, _ := os.Getwd()

		fmt.Println("Created file at path " + wd + "/" + BundleMakeFileName)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)
}

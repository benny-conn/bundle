package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/internal"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup structure for running a server with access to the official Bundle Repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := os.Getwd()
		if err != nil {
			return err
		}

		if len(args) > 0 {
			path = args[0]
		}

		valid := internal.IsValidPath(path)
		if !valid {
			return errors.New("invalid path")
		}

		if file.IsBundleInitialized(path) {
			return errors.New("there already exists a bundle.yml at this location")
		}

		if !isPluginDirectory(path) {
			return errors.New("there is no plugin directory in your current directory")
		}

		err = file.Initialize(path)
		if err != nil {
			return err
		}
		fmt.Println("Created file at path " + path + "/" + file.BuFileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

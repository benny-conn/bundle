package cli

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/intfile"
	"github.com/bennycio/bundle/cli/uploader"
	"github.com/bennycio/bundle/internal"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your plugin as specified in bundle-make.yml to the official Bundle Repository",
	Long: `Will upload the jar specified under JarPath into the official Bundle Repository, allowing public access
	to your plugin. Version must be unique per upload and name must be unique globally for the initial upload`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if !internal.IsValidPath(args[0]) {
			log.Fatal(errors.New("invalid path").Error())
		}

		path := args[0]

		user := credentialsPrompt()

		plugin := &api.Plugin{}

		isReadme := strings.HasSuffix(path, "README.md")

		if isReadme {

			plugin = pluginInfoPrompt()

		} else {
			result, err := intfile.ParsePluginYml(path)

			if err != nil {
				return err
			}

			plugin.Name = result.Name
			plugin.Version = result.Version
			plugin.Description = result.Description
		}

		fmt.Printf("Uploading to Bundle Repository From: %s\n", path)

		upl := uploader.New(user, path, plugin.Name, plugin.Version).WithReadme(isReadme)

		err := upl.Upload()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func pluginInfoPrompt() *api.Plugin {
	fmt.Println("Enter plugin name: ")
	var pluginName string
	fmt.Scanln(&pluginName)

	plugin := &api.Plugin{
		Name:    pluginName,
		Version: "README",
	}

	return plugin
}

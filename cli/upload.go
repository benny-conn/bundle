package cli

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/intfile"
	"github.com/bennycio/bundle/cli/term"
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
			plugin.Category = api.Category(result.Category)

		}

		term.Print(fmt.Sprintf("Uploading to Bundle Repository From: %s\n", path))

		fi, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fi.Close()

		upl := uploader.New(user, fi, plugin).WithReadme(isReadme)

		err = upl.Upload()
		if err != nil {
			return err
		}
		term.Println("Successfully uploaded! :)")

		if !isReadme {
			term.Println("Would you like to upload a README as well? [Y/n]")

			var rdmeToo string
			fmt.Scanln(&rdmeToo)

			if strings.EqualFold(rdmeToo, "y") || strings.EqualFold(rdmeToo, "yes") {
				term.Println("Please specify a path to your readme file or press enter to scan for readme in close directories.")
				rd := bufio.NewReader(os.Stdin)
				p, err := rd.ReadString(byte('\n'))
				if err != nil {
					return err
				}
				p = strings.TrimSpace(strings.Trim(p, "\n"))
				var readme *os.File
				if p == "" {
					wlk := uploader.NewFileWalker("README.md", plugin.Name)
					readme, err = wlk.Walk()
					if err != nil {
						return err
					}
				} else {
					readme, err = os.Open(path)
					if err != nil {
						return err
					}
				}
				rdmeUpl := uploader.New(user, readme, plugin).WithReadme(true)
				err = rdmeUpl.Upload()
				if err != nil {
					return err
				}
				term.Println("Successfully uploaded! :)")
			}
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

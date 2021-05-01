package cli

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your plugin as specified in bundle-make.yml to the official Bundle Repository",
	Long: `Will upload the jar specified under JarPath into the official Bundle Repository, allowing public access
	to your plugin. Version must be unique per upload and name must be unique globally for the initial upload`,
	Run: func(cmd *cobra.Command, args []string) {

		if !bundle.IsValidPath(args[0]) {
			log.Fatal(errors.New("invalid path").Error())
		}

		path := args[0]

		user := credentialsPrompt()

		plugin := &api.Plugin{}

		isReadme := strings.HasSuffix(path, "README.md")

		if isReadme {

			plugin = pluginInfoPrompt()

		} else {
			reader, err := zip.OpenReader(path)

			if err != nil {
				log.Fatal(err)
			}

			result := &PluginYML{}

			for _, file := range reader.File {
				if strings.HasSuffix(file.Name, "plugin.yml") {
					rc, err := file.Open()
					if err != nil {
						log.Fatal(err)
					}
					buf := bytes.Buffer{}
					buf.ReadFrom(rc)
					err = yaml.Unmarshal(buf.Bytes(), result)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			reader.Close()

			plugin.Name = result.Name
			plugin.Version = result.Version

		}

		fmt.Printf("Uploading to Bundle Repository From: %s\n", path)

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if isReadme {
			err = wrapper.UploadReadmeApi(user, plugin.Name, file)
		} else {
			err = wrapper.UploadPluginApi(user, plugin.Name, plugin.Version, file)
		}
		if err != nil {
			panic(err)
		}

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

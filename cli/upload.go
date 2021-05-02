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
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate"
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
			plugin.Description = result.Description

		}

		fmt.Printf("Uploading to Bundle Repository From: %s\n", path)

		gservice := gate.NewGateService("localhost", "8020")

		if isReadme {
			file, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			dbPlugin, err := gservice.GetPlugin(plugin)
			if err != nil {
				panic(err)
			}
			readme := &api.Readme{
				Plugin: dbPlugin.Id,
				Text:   string(file),
			}
			err = gservice.InsertReadme(user, readme)
		} else {
			file, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			err = gservice.UploadPlugin(user, plugin, file)
			if err != nil {
				panic(err)
			}
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

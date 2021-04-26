package cli

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	bundle "github.com/bennycio/bundle/internal"
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

		plugin := bundle.Plugin{}

		isReadme := strings.HasSuffix(path, "README.md")

		if isReadme {

			plugin = pluginInfoPrompt()

		} else {
			reader, err := zip.OpenReader(path)

			if err != nil {
				log.Fatal(err)
			}

			result := &bundle.PluginYML{}

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
			plugin.Author = user.Username
		}

		fmt.Printf("Uploading to Bundle Repository From: %s\n", path)

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}

		resp, err := uploadToRepo(file, plugin, user)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resp.Status))
		fmt.Println(string(respBody))
	},
}

func uploadToRepo(file io.Reader, plugin bundle.Plugin, user bundle.User) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/bundle", file)
	if err != nil {
		return nil, err
	}

	userJSON, err := json.Marshal(user)

	if err != nil {
		log.Fatal(err)
	}

	pluginJSON, err := json.Marshal(plugin)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Plugin", string(pluginJSON))
	req.Header.Add("User", string(userJSON))

	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func pluginInfoPrompt() bundle.Plugin {
	fmt.Println("Enter plugin name: ")
	var pluginName string
	fmt.Scanln(&pluginName)

	plugin := bundle.Plugin{
		Name:    pluginName,
		Version: "README",
	}

	return plugin
}

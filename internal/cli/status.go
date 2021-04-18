package cli

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check which plugins have updates.",
	Run: func(cmd *cobra.Command, args []string) {
		m := getBundledPlugins()

		pluginsToUpdate := make(map[string]string)

		for k, v := range m {
			var updatedVersion string

			if v != "latest" {
				updatedVersion = v
			} else {
				plugin, err := getPlugin(k)
				if err != nil {
					panic(err)
				}

				updatedVersion = plugin.Version
			}

			fp := filepath.Join("plugins", k+".jar")

			reader, err := zip.OpenReader(fp)

			if err != nil {
				panic(err)
			}

			for _, file := range reader.File {
				if strings.HasSuffix(file.Name, "plugin.yml") {
					yml := &PluginYML{}
					rc, err := file.Open()
					if err != nil {
						panic(err)
					}
					buf := bytes.Buffer{}
					buf.ReadFrom(rc)
					yaml.Unmarshal(buf.Bytes(), yml)

					if yml.Version != updatedVersion {
						pluginsToUpdate[k] = updatedVersion
					}
				}
			}

		}

		if len(pluginsToUpdate) != 0 {
			fmt.Println("Plugins To Update:")
			for k, v := range pluginsToUpdate {
				fmt.Println(k, " -> ", v)
			}
		} else {
			fmt.Println("All plugins are up to date :)")
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

package cli

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/intfile"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check which plugins have updates.",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := intfile.GetBundleFilePlugins(buFilePath)
		if err != nil {
			panic(err)
		}

		pluginsToUpdate := make(map[string]string)

		var wg sync.WaitGroup

		wg.Add(len(m))

		for k, v := range m {
			go func(pluginName string, bundleVersion string) {
				defer wg.Done()

				req := &api.Plugin{
					Name: pluginName,
				}
				gs := gate.NewGateService("localhost", "8020")

				plugin, err := gs.GetPlugin(req)
				if err != nil {
					panic(err)
				}

				latestVersion := plugin.Version

				fp := filepath.Join(buFilePath, "plugins", pluginName+".jar")

				res, err := intfile.ParsePluginYml(fp)
				if err != nil {
					fmt.Printf("error occurred: %s\n", res)
				}
				if res.Version != latestVersion {
					pluginsToUpdate[pluginName] = latestVersion
				}
			}(k, v)
		}
		wg.Wait()
		if len(pluginsToUpdate) > 0 {
			fmt.Println("Plugins To Update:")
			for k, v := range pluginsToUpdate {
				fmt.Println(k, " -> ", v)
			}
			fmt.Println(`Use "bundle install" to update :)`)
		} else {
			fmt.Println("All plugins are up to date :)")
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

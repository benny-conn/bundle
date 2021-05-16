package cli

import (
	"fmt"
	"os"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := intfile.GetBundle("")
		if err != nil {
			return err
		}

		pls := bundle.Plugins

		pluginsToUpdate := make(map[string]string)

		var wg sync.WaitGroup

		wg.Add(len(pls))

		for k, v := range pls {
			go func(pluginName string, bundleVersion string) {
				defer wg.Done()

				req := &api.Plugin{
					Name: pluginName,
				}
				gs := gate.NewGateService("localhost", "8020")

				plugin, err := gs.GetPlugin(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error occurred: %s\n", err.Error())
				}

				latestVersion := plugin.Version

				fp := filepath.Join("plugins", pluginName+".jar")

				res, err := intfile.ParsePluginYml(fp)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error occurred: %s\n", err.Error())
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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

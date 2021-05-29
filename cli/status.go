package cli

import (
	"fmt"
	"os"
	"sync"

	"github.com/alexeyco/simpletable"
	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check which plugins have updates.",
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := file.GetBundle("")
		if err != nil {
			return err
		}

		printStatus(bundle.Plugins, nil)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func printStatus(pls map[string]string, conn *ftp.ServerConn) {

	pluginsToUpdate := sync.Map{}

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
				return
			}

			latestVersion := plugin.Version

			if conn == nil {
				if res, err := file.GetPluginYml(pluginName, nil); err == nil {
					pls[pluginName] = res.Version
					if res.Version != latestVersion {
						pluginsToUpdate.Store(pluginName, latestVersion)
					}
				} else {
					pls[pluginName] = "Not Installed"
					pluginsToUpdate.Store(pluginName, latestVersion)
				}

			} else {
				if res, err := file.GetPluginYml(pluginName, conn); err == nil {
					pls[pluginName] = res.Version
					if res.Version != latestVersion {
						pluginsToUpdate.Store(pluginName, latestVersion)
					}
				} else {
					pls[pluginName] = "Not Installed"
					pluginsToUpdate.Store(pluginName, latestVersion)
				}
			}

		}(k, v)
	}
	wg.Wait()

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Plugin"},
			{Text: "Current"},
			{Text: "Updated"},
		},
	}

	pluginsToUpdate.Range(func(key, value interface{}) bool {
		r := []*simpletable.Cell{
			{Text: key.(string)},
			{Text: pls[key.(string)]},
			{Text: value.(string)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
		return true
	})

	table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())
	if conn == nil {
		term.Println(`Use "bundle install" to update your plugins`)
	} else {
		term.Println(`Use "install" to update your plugins`)
	}
}

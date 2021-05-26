package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexeyco/simpletable"
	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal/gate"
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

		findStatuses(bundle.Plugins)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func findStatuses(pls map[string]string) {

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

			fp := filepath.Join("plugins", pluginName+".jar")

			plfile, err := os.Open(fp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error occurred: %s\n", err.Error())
				return
			}
			defer plfile.Close()

			res, err := file.ParsePluginYml(plfile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error occurred: %s\n", err.Error())
				return
			}
			pls[pluginName] = res.Version
			if res.Version != latestVersion {
				pluginsToUpdate.Store(pluginName, latestVersion)
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
	term.Println(`Use "bundle install" to update your plugins`)
}

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var SpecifiedVersion string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"update", "get"},
	Short:   "Install plugins for your Bundle",
	Long: `Install plugins from the official Bundle Repository to your Bundle. If no plugins are
	specified, all plugins listed in bundle.yml will be downloaded. Any arguments to this command
	will be interpreted as plugins to fetch from the Bundle Repository, add to your bundle.yml, and 
	download to your plugins folder`,
	Run: func(cmd *cobra.Command, args []string) {

		bundlePlugins, err := getBundleFilePlugins()
		if err != nil {
			panic(err)
		}
		if bundlePlugins == nil {
			bundlePlugins = make(map[string]string)
		}

		if args[0] != "" {
			_, err = installPlugin(args[0], SpecifiedVersion)
			if err != nil {
				panic(err)
			}
			bundlePlugins[args[0]] = SpecifiedVersion
			writePluginsToBundle(bundlePlugins)
			return
		}

		var wg sync.WaitGroup
		length := len(bundlePlugins)
		wg.Add(length)
		totalProgressBar := progressbar.Default(int64(length))
		for k, v := range bundlePlugins {
			go func(pluginName string, bundleVersion string) {
				defer wg.Done()
				defer totalProgressBar.Add(1)
				finalVersion, err := installPlugin(pluginName, bundleVersion)
				if finalVersion != v && v != "latest" {
					bundlePlugins[k] = finalVersion
				}
				if err != nil {
					panic(err)
				}
			}(k, v)
		}
		err = writePluginsToBundle(bundlePlugins)
		if err != nil {
			panic(err)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVarP(&SpecifiedVersion, "version", "v", "latest", "Specify version for installing")
}

func installPlugin(pluginName string, bundleVersion string) (string, error) {

	req := &api.Plugin{
		Name:    pluginName,
		Version: bundleVersion,
	}

	gservice := gate.NewGateService("localhost", "8020")

	if Force && req.Version != "latest" {
		plugin, err := gservice.GetPlugin(req)
		if err != nil {
			return "", err
		}
		req.Version = plugin.Version
	}

	fmt.Printf("Installing Jar %s with version %s\n", req.Name, req.Version)

	pl, err := gservice.DownloadPlugin(req)
	if err != nil {
		return "", err
	}

	fp := filepath.Join("plugins", req.Name+".jar")

	file, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	err = file.Truncate(0)
	if err != nil {
		return "", err
	}

	file.Write(pl)

	fmt.Printf("Successfully downloaded the plugin %s with version %s at file path: %s \n", pluginName, bundleVersion, file.Name())
	return req.Version, nil
}

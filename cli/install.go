package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bennycio/bundle/cli/downloader"
	"github.com/bennycio/bundle/cli/intfile"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"update", "get", "download"},
	Short:   "Install plugins for your Bundle",
	Long: `Install plugins from the official Bundle Repository to your Bundle. If no plugins are
	specified, all plugins listed in bundle.yml will be downloaded. Any arguments to this command
	will be interpreted as plugins to fetch from the Bundle Repository, add to your bundle.yml, and 
	download to your plugins folder`,
	RunE: func(cmd *cobra.Command, args []string) error {

		bundle, err := intfile.GetBundle("")
		if err != nil {
			return err
		}

		bundlePlugins := bundle.Plugins
		if bundlePlugins == nil {
			bundlePlugins = make(map[string]string)
		}

		var wg sync.WaitGroup
		var length int
		if len(args) > 0 {
			length = len(args)
		} else {
			length = len(bundlePlugins)
		}

		wg.Add(length)
		totalProgressBar := progressbar.NewOptions(length, progressbar.OptionFullWidth(), progressbar.OptionSetItsString("pls"), progressbar.OptionShowCount(), progressbar.OptionShowIts(), progressbar.OptionClearOnFinish())
		if len(args) > 0 {
			for _, v := range args {
				version := "latest"
				spl := strings.Split(v, "@")
				if len(spl) > 1 {
					version = spl[1]
				}
				go func(pluginName string, bundleVersion string) {
					defer wg.Done()
					defer totalProgressBar.Add(1)
					ver, err := downloadAndInstall(pluginName, bundleVersion)
					if err != nil {
						fmt.Printf("error occured: %s\n", err.Error())
					}
					if bundleVersion != "latest" && force {
						bundlePlugins[pluginName] = ver
					}
				}(v, version)
			}
		} else {

			for k, v := range bundlePlugins {

				go func(pluginName string, bundleVersion string) {
					defer wg.Done()
					defer totalProgressBar.Add(1)
					ver, err := downloadAndInstall(pluginName, bundleVersion)

					if err != nil {
						fmt.Printf("error occured: %s\n", err.Error())
					}
					if bundleVersion != "latest" && force {
						bundlePlugins[pluginName] = ver
					}
				}(k, v)
			}
		}

		wg.Wait()
		err = intfile.WritePluginsToBundle(bundlePlugins, "")
		if err != nil {
			return err
		}

		fmt.Println("Successfully installed plugins! :)")
		return nil
	},
}

func downloadAndInstall(pluginName string, bundleVersion string) (string, error) {
	fp := filepath.Join("plugins", pluginName+".jar")
	latest := strings.EqualFold(bundleVersion, "latest")
	dl := downloader.New(pluginName, bundleVersion).WithLocation(fp).WithLatest(latest)
	bs, err := dl.Download()
	if err != nil {
		return "", err
	}
	err = dl.Install(bs)
	if err != nil {
		return "", err
	}
	return dl.Plugin.Version, nil
}

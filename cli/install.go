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
	Aliases: []string{"update", "get"},
	Short:   "Install plugins for your Bundle",
	Long: `Install plugins from the official Bundle Repository to your Bundle. If no plugins are
	specified, all plugins listed in bundle.yml will be downloaded. Any arguments to this command
	will be interpreted as plugins to fetch from the Bundle Repository, add to your bundle.yml, and 
	download to your plugins folder`,
	RunE: func(cmd *cobra.Command, args []string) error {

		bundlePlugins, err := intfile.GetBundleFilePlugins(buFilePath)
		if err != nil {
			return err
		}
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
		totalProgressBar := progressbar.Default(int64(length))
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
					if strings.EqualFold(bundleVersion, "latest") && force {
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
					if strings.EqualFold(bundleVersion, "latest") && force {
						bundlePlugins[pluginName] = ver
					}
				}(k, v)
			}
		}

		err = intfile.WritePluginsToBundle(bundlePlugins, buFilePath)
		if err != nil {
			return err
		}
		wg.Wait()
		return nil
	},
}

func downloadAndInstall(pluginName string, bundleVersion string) (string, error) {
	fp := filepath.Join("plugins", pluginName+".jar")
	latest := strings.EqualFold(bundleVersion, "latest")
	dl := downloader.New(pluginName, bundleVersion).WithForce(force).WithLocation(fp).WithLatest(latest)
	bs, err := dl.Download()
	if err != nil {
		return "", err
	}
	err = dl.Install(bs)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return dl.Plugin.Version, nil
}

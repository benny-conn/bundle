package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type downloadedPlugin struct {
	Plugin *api.Plugin
	Data   []byte
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force installation without approval of changes and forcibly updates versions")
}

var force bool

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

		bundle, err := file.GetBundle("")
		if err != nil {
			return err
		}

		bundlePlugins := bundle.Plugins
		if bundlePlugins == nil {
			bundlePlugins = make(map[string]string)
		}

		if len(args) > 0 {

			plsToInst := map[string]string{}
			for _, v := range args {
				spl := strings.Split(v, "@")
				if len(spl) < 2 {
					plsToInst[spl[0]] = "latest"
				} else {
					plsToInst[spl[0]] = spl[1]
				}
			}
			if err := downloadAndInstall(plsToInst); err != nil {
				return err
			}
		} else {
			if err := downloadAndInstall(bundlePlugins); err != nil {
				return err
			}
		}

		term.Println(Green("Successfully installed plugins! :)").Bold())
		return nil
	},
}

// func downloadAndInstall(pluginName string, bundleVersion string) (string, error) {

// 	fp := filepath.Join("plugins", pluginName+".jar")
// 	latest := strings.EqualFold(bundleVersion, "latest")
// 	dl := downloader.New(pluginName, bundleVersion).WithLocation(fp).WithLatest(latest)
// 	bs, err := dl.Download()
// 	if err != nil {
// 		return "", err
// 	}
// 	err = dl.Install(bs)
// 	if err != nil {
// 		return "", err
// 	}
// 	return dl.Plugin.Version, nil
// }

func changesSinceCurrent(pluginId, pluginName, desiredVersion, currentVersion string) ([]string, error) {
	gs := gate.NewGateService("localhost", "8020")
	ch := &api.Changelog{PluginId: pluginId}

	resp, err := gs.GetChangelogs(ch)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s %s\n", Blue("Changes Since Last Update"), Blue(pluginName).Bold())

	versionsSinceUpdate := []string{"latest"}
	for _, v := range resp.Changelogs {
		if versionGreaterThan(v.Version, currentVersion) {
			if versionGreaterThan(desiredVersion, v.Version) || desiredVersion == v.Version {
				versionsSinceUpdate = append(versionsSinceUpdate, v.Version)
				fmt.Println(Yellow(v.Version).Bold())
				fmt.Println(Green("Added: ").Bold())
				for _, v := range v.Added {
					fmt.Printf("  - %s\n", Green(v))
				}
				fmt.Println(Red("Removed: ").Bold())
				for _, v := range v.Removed {
					fmt.Printf("  - %s\n", Red(v))
				}
				fmt.Println(Blue("Updated: ").Bold())
				for _, v := range v.Updated {
					fmt.Printf("  - %s\n", Blue(v))
				}
			}
		}
	}
	return versionsSinceUpdate, nil
}

func downloadAndInstall(plugins map[string]string) error {
	gs := gate.NewGateService("localhost", "8020")
	installQueue := make(chan downloadedPlugin)
	mu := &sync.Mutex{}
	i := 1
	for k, v := range plugins {
		go func(index int, pluginName, version string) {
			if len(plugins) <= index {
				defer close(installQueue)
			}
			pl := &api.Plugin{Name: pluginName}
			dbpl, err := gs.GetPlugin(pl)
			if err != nil {
				fmt.Printf("error occurred: %s", err.Error())
				return
			}
			pl.Id = dbpl.Id
			pl.Name = dbpl.Name

			if strings.EqualFold(version, "latest") || version == "" {
				pl.Version = dbpl.Version
			} else {
				pl.Version = version
			}

			plfile, err := os.Open(fmt.Sprintf("plugins/%s.jar", pl.Name))
			if err == nil {
				defer plfile.Close()

				plyml, err := file.ParsePluginYml(plfile)
				if err != nil {
					fmt.Printf("error occurred: %s", err.Error())
					return
				}

				downloadedVer := plyml.Version

				if downloadedVer == dbpl.Version || downloadedVer == pl.Version {
					return
				}

				mu.Lock()
				func() {
					defer mu.Unlock()
					missedVers, err := changesSinceCurrent(pl.Id, pl.Name, pl.Version, plyml.Version)
					if err != nil {
						fmt.Printf("error occurred: %s", err.Error())
						return
					}
					term.Println(fmt.Sprintf("Which version would you like to update to for the plugin: %s (%d/%d)?\nType 'latest' for the latest version", pl.Name, index, len(plugins)))
					resVer := prompt.Choose(">> ", missedVers)
					if !strings.EqualFold(resVer, "latest") {
						pl.Version = resVer
					}
				}()
			}

			bs, err := gs.DownloadPlugin(pl)
			if err != nil {
				fmt.Printf("error occurred: %s", err.Error())
				return
			}
			installQueue <- downloadedPlugin{Plugin: pl, Data: bs}
		}(i, k, v)
		i += 1
	}
	for v := range installQueue {
		func() {

			pb := progressbar.NewOptions(
				len(v.Data),
				progressbar.OptionClearOnFinish(),
				progressbar.OptionSetDescription(fmt.Sprintf("Installing %s - %s", v.Plugin.Name, v.Plugin.Version)),
				progressbar.OptionShowBytes(true),
				progressbar.OptionShowCount(),
				progressbar.OptionSetItsString("bytes"),
			)
			pr, pw := io.Pipe()

			go func() {
				defer pw.Close()

				writer := io.MultiWriter(pb, pw)
				_, err := writer.Write(v.Data)
				if err != nil {
					fmt.Printf("error occurred: %s", err.Error())
					return
				}
			}()

			fp := filepath.Join("plugins", v.Plugin.Name+".jar")
			os.Remove(fp)
			fi, err := os.Create(fp)
			if err != nil {
				fmt.Printf("error occurred: %s", err.Error())
				return
			}
			defer fi.Close()
			_, err = io.Copy(fi, pr)
			if err != nil {
				fmt.Printf("error occurred: %s", err.Error())
				return
			}
		}()
	}
	return nil
}

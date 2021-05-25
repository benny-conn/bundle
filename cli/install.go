package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/downloader"
	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

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

		gs := gate.NewGateService("localhost", "8020")

		if len(args) > 0 {
			for _, v := range args {
				version := "latest"
				spl := strings.Split(v, "@")
				if len(spl) > 1 {
					version = spl[1]
				}
				ver, err := downloadAndInstall(v, version)
				if err != nil {
					fmt.Printf("error occured: %s\n", err.Error())
				} else {
					if version != "latest" && force {
						bundlePlugins[v] = ver
					}
				}
			}
		} else {
			i := 1
			for k, v := range bundlePlugins {

				fp := filepath.Join("plugins", k+".jar")
				yml, err := file.ParsePluginYml(fp)

				if err == nil {
					dbpl, err := gs.GetPlugin(&api.Plugin{Name: yml.Name})
					if err != nil {
						return err
					}
					err = changesSinceCurrent(dbpl.Id, dbpl.Name, yml.Version)
					if err != nil {
						return err
					}
					if !force {
						term.Println(fmt.Sprintf("Update Plugin %s (%d/%d)? [Y/n]", dbpl.Name, i, len(bundlePlugins)))
						cont := prompt.Input(">> ", yesOrNoCompleter)
						if !strings.EqualFold(cont, "y") && !strings.EqualFold(cont, "yes") {
							continue
						}
					}
				}

				ver, err := downloadAndInstall(k, v)

				if err != nil {
					fmt.Printf("error occured: %s\n", err.Error())
				} else {
					if v != "latest" && force {
						bundlePlugins[k] = ver
					}
				}
				i += 1
			}
		}
		err = file.WritePluginsToBundle(bundlePlugins, "")
		if err != nil {
			return err
		}

		term.Println(Green("Successfully installed plugins! :)").Bold())
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

func changesSinceCurrent(pluginId, pluginName, currentVersion string) error {
	gs := gate.NewGateService("localhost", "8020")
	ch := &api.Changelog{PluginId: pluginId}

	resp, err := gs.GetChangelogs(ch)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s\n", Blue("Changes Since Last Update"), Blue(pluginName).Bold())

	for _, v := range resp.Changelogs {
		if versionGreaterThan(v.Version, currentVersion) {
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
	return nil
}

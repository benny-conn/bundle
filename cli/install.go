package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/cli/file"
	"github.com/bennycio/bundle/cli/logger"
	"github.com/bennycio/bundle/cli/term"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/c-bata/go-prompt"
	"github.com/jlaffaye/ftp"
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

		user, err := getCurrentUser()
		if err != nil {
			return err
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
			if err := downloadAndInstall(plsToInst, user, nil); err != nil {
				return err
			}
		} else {
			if err := downloadAndInstall(bundlePlugins, user, nil); err != nil {
				return err
			}
		}

		term.Println(Green("Successfully installed plugins! :)").Bold())
		return nil
	},
}

func changesSinceCurrent(pluginId, pluginName, desiredVersion, currentVersion string) ([]string, error) {
	gs := gate.NewGateService("localhost", "8020")
	ch := &api.Changelog{PluginId: pluginId}

	resp, err := gs.GetChangelogs(ch)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s %s\n", Blue("Changes Since Last Update"), Blue(pluginName).Bold())

	versionsSinceUpdate := []string{}
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

func downloadAndInstall(plugins map[string]string, user *api.User, conn *ftp.ServerConn) error {
	gs := gate.NewGateService("localhost", "8020")
	installQueue := make(chan downloadedPlugin)
	mu := &sync.Mutex{}
	left := int64(len(plugins))
	i := 1
	for k, v := range plugins {
		go func(index int, pluginName, version string) {
			defer atomic.AddInt64(&left, -1)
			pl := &api.Plugin{Name: pluginName}
			dbpl, err := gs.GetPlugin(pl)
			if err != nil {
				logger.ErrLog.Print(err.Error())
			}
			pl.Id = dbpl.Id
			pl.Name = dbpl.Name

			if strings.EqualFold(version, "latest") || version == "" {
				pl.Version = dbpl.Version
			} else {
				pl.Version = version
			}

			if conn == nil {
				if plyml, err := file.GetPluginYml(pluginName, nil); err == nil {
					if plyml.Version == dbpl.Version || plyml.Version == pl.Version {
						return
					}
					mu.Lock()
					func() {
						defer mu.Unlock()
						missedVers, err := changesSinceCurrent(pl.Id, pl.Name, pl.Version, plyml.Version)
						if err != nil {
							logger.ErrLog.Print(err.Error())
						}
						term.Println(fmt.Sprintf("Which version would you like to update to for the plugin: %s (%d/%d)?\nPress enter for the latest version", pl.Name, index, len(plugins)))
						resVer := prompt.Choose(">> ", missedVers)
						if resVer != "" {
							pl.Version = resVer
						}
					}()
				}
			} else {
				if plyml, err := file.GetPluginYml(pluginName, conn); err == nil {
					if plyml.Version == dbpl.Version || plyml.Version == pl.Version {
						return
					}
					mu.Lock()
					func() {
						defer mu.Unlock()
						missedVers, err := changesSinceCurrent(pl.Id, pl.Name, pl.Version, plyml.Version)
						if err != nil {
							logger.ErrLog.Print(err.Error())
						}
						term.Println(fmt.Sprintf("Which version would you like to update to for the plugin: %s (%d/%d)?\nPress enter for the latest version", pl.Name, index, len(plugins)))
						resVer := prompt.Choose(">> ", missedVers)
						if resVer != "" {
							pl.Version = resVer
						}
					}()
				}
			}

			bs, err := gs.DownloadPlugin(pl, user)
			if err != nil {
				logger.ErrLog.Print(err.Error())
				return
			}
			installQueue <- downloadedPlugin{Plugin: pl, Data: bs}
		}(i, k, v)
		i += 1
	}
	for {
		select {
		case v := <-installQueue:
			func() {
				pb := progressbar.DefaultBytes(int64(len(v.Data)), fmt.Sprintf("Installing %s - %s", v.Plugin.Name, v.Plugin.Version))
				pr, pw := io.Pipe()

				go func() {
					defer pw.Close()

					writer := io.MultiWriter(pb, pw)
					_, err := writer.Write(v.Data)
					if err != nil {
						logger.ErrLog.Print(err.Error())
					}
				}()

				fp := filepath.Join("plugins", v.Plugin.Name+".jar")
				if conn == nil {
					os.Remove(fp)
					fi, err := os.Create(fp)
					if err != nil {
						logger.ErrLog.Print(err.Error())
					}
					defer fi.Close()
					_, err = io.Copy(fi, pr)
					if err != nil {
						logger.ErrLog.Print(err.Error())
					}
				} else {
					err := conn.Stor(fp, pr)
					if err != nil {
						logger.ErrLog.Print(err.Error())
					}
				}
			}()
		default:
		}
		if left < 1 {
			break
		}
	}

	return nil
}

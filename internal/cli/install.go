package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/bennycio/bundle/pkg"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

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
		var wg sync.WaitGroup
		m, err := getBundleFilePlugins()
		if err != nil {
			panic(err)
		}
		length := len(m)
		wg.Add(length)
		totalProgressBar := progressbar.Default(int64(length))
		for k, v := range m {
			go func(key string, value string) {
				defer wg.Done()
				fmt.Printf("Installing Jar %s with version %s\n", key, value)

				u, err := url.Parse("http://localhost:8080/bundle")
				if err != nil {
					panic(err)
				}
				q := u.Query()
				q.Set("name", key)

				version := value
				if Force {
					plugin, err := pkg.GetPlugin(key)
					if err != nil {
						panic(err)
					}
					version = plugin.Version
				}
				q.Set("version", version)
				u.RawQuery = q.Encode()

				resp, err := http.Get(u.String())

				if err != nil {
					panic(err)
				}

				defer resp.Body.Close()

				fp := filepath.Join("plugins", key+".jar")

				file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

				if err != nil {
					panic(err)
				}

				io.Copy(file, resp.Body)
				fmt.Printf("Successfully downloaded the plugin %s with version %s at file path: %s \n", key, value, file.Name())
				totalProgressBar.Add(1)
			}(k, v)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

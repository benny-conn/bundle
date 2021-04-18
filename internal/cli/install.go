package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

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
		m := getBundledPlugins()
		for k, v := range m {
			fmt.Printf("Installing Jar %s with version %s\n", k, v)

			u, err := url.Parse("http://localhost:8080/bundle")
			if err != nil {
				panic(err)
			}
			q := u.Query()
			q.Set("name", k)
			q.Set("version", v)
			u.RawQuery = q.Encode()

			resp, err := http.Get(u.String())

			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			fp := filepath.Join("plugins", k+".jar")

			file, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				panic(err)
			}

			io.Copy(file, resp.Body)
			fmt.Printf("Successfully downloaded the plugin %s with version %s at file path: %s \n", k, v, file.Name())
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

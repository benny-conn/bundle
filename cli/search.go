package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/spf13/cobra"
)

var page int

var count int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for plugins to download",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no search parameters specified")
		}

		params := strings.Join(args, " ")

		gs := gate.NewGateService("localhost", "8020")

		results, err := gs.PaginatePlugins(&api.PaginatePluginsRequest{
			Page:   int32(page),
			Count:  int32(count),
			Search: params,
		})

		if err != nil {
			return err
		}

		for i, v := range results {
			fmt.Printf("%v. %s: \nDescription: %s\nAuthor: %s\n\n", i, v.Name, v.Description, v.Author.Username)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().IntVarP(&page, "page", "p", 1, "Select a page for search results")
	searchCmd.Flags().IntVarP(&count, "count", "c", 15, "Select how many results per page")

}

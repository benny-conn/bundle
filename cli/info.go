package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/gate"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info on a specific plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no search parameters specified")
		}
		pl := args[0]

		gs := gate.NewGateService("localhost", "8020")

		result, err := gs.GetPlugin(&api.Plugin{Name: pl})

		if err != nil {
			return err
		}

		fmt.Println(Blue("|| -- Plugin Info ------------ ||").Bold())

		fmt.Printf("Name: %s\n", result.Name)
		fmt.Printf("Author: %s\n", result.Author.Username)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("Current Version: %s\n", result.Version)

		fmt.Println("Would you like to see recent changes? [Y/n]")

		cont := prompt.Input(">> ", yesOrNoCompleter)
		if strings.EqualFold(cont, "y") || strings.EqualFold(cont, "yes") {
			ch, err := gs.GetChangelog(&api.Changelog{PluginId: result.Id, Version: result.Version})
			if err != nil {
				return err
			}
			fmt.Println(Green("Added: ").Bold())
			for _, v := range ch.Added {
				fmt.Printf("  - %s\n", Green(v))
			}
			fmt.Println(Red("Removed: ").Bold())
			for _, v := range ch.Removed {
				fmt.Printf("  - %s\n", Red(v))
			}
			fmt.Println(Blue("Updated: ").Bold())
			for _, v := range ch.Updated {
				fmt.Printf("  - %s\n", Blue(v))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

package cli

import (
	"github.com/spf13/cobra"
)

var interactive bool

// initCmd represents the init command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup structure for running a server with access to the official Bundle Repository",
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Edit configuration interactively")
}

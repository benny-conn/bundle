package cli

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check which plugins have updates.",
	Run: func(cmd *cobra.Command, args []string) {
		// m := getBundledPlugins()

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

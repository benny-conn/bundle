package cli

import (
	"github.com/bennycio/bundle/cli/term"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print path to config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		term.Println(configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

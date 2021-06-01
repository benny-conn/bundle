package cli

import (
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for Bundle MC, allowing you to upload plugins to the official repository",
	RunE: func(cmd *cobra.Command, args []string) error {

		_, err := newUser()

		if err != nil {
			return err
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

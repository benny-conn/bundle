package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for Bundle MC, allowing you to upload plugins to the official repository",
	Long:  "Sign up for Bundle MC and gain upload access to the official repository. Use flags \"-u\" \"-e\" and \"-p\" to specify username, email, and password ",
	RunE: func(cmd *cobra.Command, args []string) error {

		user := credentialsPrompt()
		creds := map[string]string{
			"username": user.Username,
			"password": user.Password,
		}
		viper.Set("Credentials", creds)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

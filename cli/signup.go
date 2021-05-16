package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for Bundle MC, allowing you to upload plugins to the official repository",
	Long:  "Sign up for Bundle MC and gain upload access to the official repository. Use flags \"-u\" \"-e\" and \"-p\" to specify username, email, and password ",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := openBrowser("https://bundlemc.io/signup")
		if err != nil {
			return err
		}
		fmt.Println("Opened signup in new browser window")

		return nil

	},
}

func init() {
	rootCmd.AddCommand(signupCmd)
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}

	return nil
}

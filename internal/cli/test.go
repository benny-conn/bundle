package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "TEST",
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if path == "" {
			log.Fatal(":(")
		}

		bytes, err := ioutil.ReadFile(path)

		if err != nil {
			panic(err)
		}

		result := http.DetectContentType(bytes)

		fmt.Println("Result is: " + result)

	},
}

func init() {
	rootCmd.AddCommand(testCmd)

}

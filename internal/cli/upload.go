package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload your plugin as specified in bundle-make.yml to the official Bundle Repository",
	Long: `Will upload the jar specified under JarPath into the official Bundle Repository, allowing public access
	to your plugin. Version must be unique per upload and name must be unique globally for the initial upload`,
	Run: func(cmd *cobra.Command, args []string) {
		// if path == "" || !strings.HasSuffix(path, ".jar") {
		// 	log.Fatal("Please specify a valid path to the jar with -p ")
		// }
		fb, err := getBundleFile(true)

		if err != nil {
			panic(err)
		}
		result := BundleMakeFile{}

		err = yaml.Unmarshal(fb, &result)

		if err != nil {
			panic(err)
		}

		path := result.JarPath
		finalName := result.Name
		version := result.Version

		user := credentialsPrompt()

		userAsJSON, err := json.Marshal(user)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Uploading to Bundle Repository From: %s\n", path)

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		// bytes, err := io.ReadAll(file)
		// if err != nil {
		// 	panic(err)
		// }

		// fileType := http.DetectContentType(bytes)

		// if fileType != REQUIRED_FILE_TYPE {
		// 	log.Fatal("File format does not meet requirements. Uploaded file must be the built Jar of your plugin")
		// }

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8070/bundle", file)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Project-Version", version)
		req.Header.Add("Project-Name", finalName)
		req.Header.Add("User", string(userAsJSON))

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resp.Status))
		fmt.Println(string(respBody))
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

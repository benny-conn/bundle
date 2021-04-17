package cli

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	bundle "github.com/bennycio/bundle/internal"
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

		if !bundle.IsValidPath(args[0]) {
			log.Fatal(errors.New("invalid path").Error())
		}

		path := args[0]

		reader, err := zip.OpenReader(path)

		if err != nil {
			log.Fatal(err)
		}

		result := &PluginYML{}

		for _, file := range reader.File {
			if strings.HasSuffix(file.Name, "plugin.yml") {
				rc, err := file.Open()
				if err != nil {
					log.Fatal(err)
				}
				buf := bytes.Buffer{}
				buf.ReadFrom(rc)
				err = yaml.Unmarshal(buf.Bytes(), result)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		reader.Close()

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

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/bundle", file)
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

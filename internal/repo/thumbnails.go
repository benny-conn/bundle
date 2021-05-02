package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bennycio/bundle/api"
	"github.com/spf13/viper"
)

func thumbnailsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		pluginJSON := r.Header.Get("Resource")

		plugin := &api.Plugin{}
		json.Unmarshal([]byte(pluginJSON), plugin)

		loc, err := uploadThumbnailToRepo(plugin, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Make gate handle db stuff
		// gs := gate.NewGateService("", "")
		// err := cleanPlugin(plugin)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// plugin.Thumbnail = fmt.Sprintf("https://bundle-repository.s3-us-east-1.amazonaws.com/%v/%v/THUMBNAIL.webp", plugin.Author, plugin.Id)

		// err = gs.UpdatePlugin(plugin)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadGateway)
		// 	return
		// }

		fmt.Println("Successfully uploaded to " + loc)
	}
}

func uploadThumbnailToRepo(plugin *api.Plugin, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(plugin.Author.Id, plugin.Id, "THUMBNAIL.webp")

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String("bundle-repository"),
		Key:    aws.String(fp),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

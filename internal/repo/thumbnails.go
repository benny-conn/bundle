package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bennycio/bundle/api"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/wrapper"
	"github.com/spf13/viper"
)

func thumbnailsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")

		req := &api.Plugin{
			Name: name,
		}

		plugin, err := wrapper.GetPluginApi(req)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pl, err := downloadThumbnailFromRepo(plugin.Name, plugin.Author)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write(pl)
	}

	if r.Method == http.MethodPost {

		pluginJSON := r.Header.Get("Resource")

		plugin := &api.Plugin{}
		json.Unmarshal([]byte(pluginJSON), plugin)

		loc, err := uploadThumbnailToRepo(plugin.Name, plugin.Author, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to " + loc)
	}
}

func uploadThumbnailToRepo(name string, author string, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(author, name, "THUMBNAIL.png")

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String("bundle-repository"),
		Key:    aws.String(fp),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func downloadThumbnailFromRepo(name string, author string) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(author, name, "THUMBNAIL.png")

	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("AWSBucket")),
		Key:    aws.String(fp),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

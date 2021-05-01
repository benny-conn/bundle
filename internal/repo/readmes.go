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

func readmesHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &api.Plugin{
			Name: r.FormValue("name"),
		}

		plugin, err := wrapper.GetPluginApi(req)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pl, err := downloadReadmeFromRepo(plugin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write(pl)
	}

	if r.Method == http.MethodPost {

		pluginJSON := r.Header.Get("Resource")

		plugin := &api.Plugin{}
		json.Unmarshal([]byte(pluginJSON), plugin)

		loc, err := uploadReadmeToRepo(plugin, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to " + loc)
	}
}

func uploadReadmeToRepo(plugin *api.Plugin, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(plugin.Author, plugin.Id, "README.md")

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

func downloadReadmeFromRepo(plugin *api.Plugin) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(plugin.Author, plugin.Id, "README.md")

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

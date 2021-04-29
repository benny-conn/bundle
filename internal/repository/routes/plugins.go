package routes

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

func PluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		err := r.ParseForm()
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := r.FormValue("name")
		version := r.FormValue("version")

		plugin, err := wrapper.GetPlugin(name)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if version != "latest" {
			plugin.Version = version
		}

		pl, err := downloadPluginFromRepo(plugin.Name, plugin.Version, plugin.Author)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write(pl)
	}

	if r.Method == http.MethodPost {

		fmt.Println("GOT HERE")

		pluginJSON := r.Header.Get("Resource")

		plugin := &api.Plugin{}
		json.Unmarshal([]byte(pluginJSON), plugin)

		err := updateOrInsertPlugin(plugin)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		loc, err := uploadPluginToRepo(plugin.Name, plugin.Version, plugin.Author, r.Body)
		if err != nil {
			bundle.WriteResponse(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to " + loc)
	}
}

func updateOrInsertPlugin(plugin *api.Plugin) error {
	dbPlugin, err := wrapper.GetPlugin(plugin.Name)

	if err == nil {
		err = wrapper.UpdatePlugin(dbPlugin.Name, plugin)
		if err != nil {
			return err
		}
	} else {

		err = wrapper.InsertPlugin(plugin)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func uploadPluginToRepo(name string, version string, author string, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fp := filepath.Join(author, name, version, name+".jar")

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

func downloadPluginFromRepo(name string, version string, author string) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fn := filepath.Join(author, name, version, name+".jar")
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("AWSBucket")),
		Key:    aws.String(fn),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

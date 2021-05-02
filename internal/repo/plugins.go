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
	"github.com/spf13/viper"
)

func pluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// we cant be usin no names no more
		req := &api.Plugin{
			Name:    r.FormValue("name"),
			Version: r.FormValue("version"),
		}

		// make gate handle all db stuff
		// gs := gate.NewGateService("", "")
		// plugin, err := gs.GetPlugin(req)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }

		// if req.Version != "latest" {
		// 	plugin.Version = req.Version
		// }

		pl, err := downloadPluginFromRepo(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write(pl)
	case http.MethodPost:
		pluginJSON := r.Header.Get("Resource")

		plugin := &api.Plugin{}

		err := json.Unmarshal([]byte(pluginJSON), plugin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		loc, err := uploadPluginToRepo(plugin, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Println("Successfully uploaded to " + loc)
	}

}

func uploadPluginToRepo(plugin *api.Plugin, body io.Reader) (string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})
	if err != nil {
		return "", err
	}

	fp := filepath.Join(plugin.Author.Id, plugin.Id, plugin.Version, plugin.Id+".jar")
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

func downloadPluginFromRepo(plugin *api.Plugin) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	fn := filepath.Join(plugin.Author.Id, plugin.Id, plugin.Version, plugin.Id+".jar")
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

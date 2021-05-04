package repo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bennycio/bundle/api"
)

func pluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req := &api.Plugin{
			Id: r.FormValue("id"),
			Author: &api.User{
				Id: r.FormValue("author"),
			},
			Version: r.FormValue("version"),
		}

		pl, err := downloadPluginFromRepo(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Write(pl)
	case http.MethodPost:

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &api.Plugin{
			Id: r.FormValue("id"),
			Author: &api.User{
				Id: r.FormValue("author"),
			},
			Version: r.FormValue("version"),
		}

		f, h, err := r.FormFile("plugin")
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if h.Size > (1024 << 20) {
			http.Error(w, "file too large", http.StatusBadRequest)
			return
		}

		loc, err := uploadPluginToRepo(req, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// UPDATE PLUGIN DB WITH THE LOCATION TO THE FILE POSSIBLY DELETE IF FAILURE :)(

		fmt.Println("Successfully uploaded to " + loc)
	}

}

func uploadPluginToRepo(plugin *api.Plugin, file io.Reader) (string, error) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return "", err
	}

	fp := filepath.Join(plugin.Author.Id, plugin.Id, plugin.Version, plugin.Id+".jar")
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String("bundle-repository"),
		Key:    aws.String(fp),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func downloadPluginFromRepo(plugin *api.Plugin) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	fn := filepath.Join(plugin.Author.Id, plugin.Id, plugin.Version, plugin.Id+".jar")
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(fn),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

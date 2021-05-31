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
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/logger"
)

func pluginsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			internal.HttpError(w, err, http.StatusInternalServerError)
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
			internal.HttpError(w, err, http.StatusServiceUnavailable)
			return
		}

		logger.DebugLog.Printf("downloading %v", req)

		w.Write(pl)
	case http.MethodPost:

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			internal.HttpError(w, err, http.StatusBadRequest)
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
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		if h.Size > (1024 << 20) {
			err = fmt.Errorf("file too large")
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		loc, err := uploadPluginToRepo(req, f)
		if err != nil {
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		logger.InfoLog.Printf("uploaded plugin with id: %s to %s", r.FormValue("id"), loc)
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

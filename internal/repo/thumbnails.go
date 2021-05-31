package repo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
)

func thumbnailsHandlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		user := r.FormValue("user")
		plugin := r.FormValue("plugin")
		file, h, err := r.FormFile("thumbnail")
		if err != nil {
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		if h.Size > (1 << 20) {
			err = fmt.Errorf("file too large")
			internal.HttpError(w, err, http.StatusBadRequest)
			return
		}

		var loc string

		if plugin != "" {
			author := r.FormValue("author")
			if author == "" {
				err = fmt.Errorf("no author specified")
				internal.HttpError(w, err, http.StatusBadRequest)
				return
			}
			pl := &api.Plugin{
				Id: plugin,
				Author: &api.User{
					Id: r.FormValue("author"),
				},
			}

			loc, err = uploadPluginThumbnail(pl, file)
			if err != nil {
				internal.HttpError(w, err, http.StatusServiceUnavailable)
				return
			}
			fmt.Println("Successfully uploaded to " + loc)
		} else if user != "" {
			u := &api.User{
				Id: user,
			}
			loc, err = uploadUserThumbnail(u, file)
			if err != nil {
				internal.HttpError(w, err, http.StatusServiceUnavailable)
				return
			}
		}
		fmt.Println("Successfully uploaded to " + loc)
	}
}

func uploadPluginThumbnail(plugin *api.Plugin, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	fp := filepath.Join(plugin.Author.Id, plugin.Id, "THUMBNAIL.webp")

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(fp),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func uploadUserThumbnail(user *api.User, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	fp := filepath.Join(user.Id, "THUMBNAIL.webp")

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(fp),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

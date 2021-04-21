package storage

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	bundle "github.com/bennycio/bundle/internal"
	"github.com/spf13/viper"
)

func UploadPlugin(plugin bundle.Plugin, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})
	fp := filepath.Join(plugin.User, plugin.Plugin, plugin.Version, plugin.Plugin+".jar")

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

func DownloadPlugin(plugin bundle.Plugin) ([]byte, error) {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	buf := aws.NewWriteAtBuffer([]byte{})

	fn := filepath.Join(plugin.User, plugin.Plugin, plugin.Version, plugin.Plugin+".jar")

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

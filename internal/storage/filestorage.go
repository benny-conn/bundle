package storage

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bennycio/bundle/api"
	"github.com/spf13/viper"
)

func UploadToRepo(plugin *api.Plugin, body io.Reader) (string, error) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})
	var fp string
	if plugin.Version == "README" {
		fp = filepath.Join(plugin.Author, plugin.Name, "README.md")
	} else {
		fp = filepath.Join(plugin.Author, plugin.Name, plugin.Version, plugin.Name+".jar")
	}

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

func DownloadFromRepo(plugin *api.Plugin, opts *api.GetPluginDataRequest) error {

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(viper.GetString("AWSRegion"))})

	if opts.WithReadme {
		fn := filepath.Join(plugin.Author, plugin.Name, "README.md")
		buf := aws.NewWriteAtBuffer([]byte{})
		downloader := s3manager.NewDownloader(sess)
		_, err := downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String(viper.GetString("AWSBucket")),
			Key:    aws.String(fn),
		})
		if err != nil {
			return err
		}
		plugin.Readme = string(buf.Bytes())
	}
	if opts.WithThumbnail {
		fn := filepath.Join(plugin.Author, plugin.Name, "THUMBNAIL.png")
		buf := aws.NewWriteAtBuffer([]byte{})
		downloader := s3manager.NewDownloader(sess)
		_, err := downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String(viper.GetString("AWSBucket")),
			Key:    aws.String(fn),
		})
		if err != nil {
			return err
		}
		plugin.Thumbnail = buf.Bytes()
	}
	if opts.WithPlugin {
		fn := filepath.Join(plugin.Author, plugin.Name, plugin.Version, plugin.Name+".jar")
		buf := aws.NewWriteAtBuffer([]byte{})
		downloader := s3manager.NewDownloader(sess)
		_, err := downloader.Download(buf, &s3.GetObjectInput{
			Bucket: aws.String(viper.GetString("AWSBucket")),
			Key:    aws.String(fn),
		})
		if err != nil {
			return err
		}
		plugin.PluginData = buf.Bytes()
	}

	return nil
}

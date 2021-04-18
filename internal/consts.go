package internal

import "github.com/spf13/viper"

const (
	BundleFileName     = "bundle.yml"
	BundleMakeFileName = "bundle-make.yml"
	RequiredFileType   = "application/zip"
)

var (
	AwsS3Region = viper.GetString("AWSRegion")
	AwsS3Bucket = viper.GetString("AWSBucket")
	MongoURL    = viper.GetString("MongoURL")
)

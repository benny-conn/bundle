package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/repo"
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/repo")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {

	id := viper.GetString("AWSID")
	secret := viper.GetString("AWSSecret")

	os.Setenv("AWS_SECRET_ACCESS_KEY", secret)
	os.Setenv("AWS_ACCESS_KEY_ID", id)

	port := os.Getenv("REPO_PORT")
	mux := repo.NewRepositoryMux()

	fmt.Printf("Started Repo Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

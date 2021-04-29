package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/repository"
	"github.com/spf13/viper"
)

func init() {
	internal.InitEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/repo")
	viper.AddConfigPath("/etc/bundle/")
	viper.AddConfigPath("$HOME/.bundle")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {

	port := os.Getenv("REPO_PORT")
	mux := repository.NewRepositoryMux()

	fmt.Printf("Started Repo Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/api"
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/api")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	port := os.Getenv("API_PORT")
	mux := api.NewApiMux()

	fmt.Printf("Started Api Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
	"github.com/spf13/viper"
)

func init() {
	internal.InitEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/web")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {

	port := os.Getenv("WEB_PORT")
	mux := web.NewWebMux()

	fmt.Printf("Started Web Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

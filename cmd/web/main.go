package main

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal/web"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/bundle/")
	viper.AddConfigPath("$HOME/.bundle")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	mux := web.NewWebMux()

	fmt.Println("Started server on port 8080")

	http.ListenAndServe(":8080", mux)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
	"github.com/spf13/viper"
)

func init() {
	internal.InitConfig()
}

func main() {

	port := viper.GetInt("Port")
	mux := web.NewWebMux()

	fmt.Printf("Started server on port %d", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

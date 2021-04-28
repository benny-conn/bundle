package main

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/api"
	"github.com/spf13/viper"
)

func init() {
	internal.InitConfig()
}

func main() {
	port := viper.GetInt("Port")
	mux := api.NewApiMux()

	fmt.Printf("Started Api Server on port %d\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

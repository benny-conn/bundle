package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/api"
)

func init() {
	internal.InitEnv()
}

func main() {
	port := os.Getenv("API_PORT")
	mux := api.NewApiMux()

	fmt.Printf("Started Api Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

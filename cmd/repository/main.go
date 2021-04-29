package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
)

func init() {
	internal.InitConfig()
}

func main() {

	port := os.Getenv("REPO_PORT")
	mux := web.NewWebMux()

	fmt.Printf("Started Repo Server on port %v", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

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

	port := os.Getenv("WEB_PORT")
	mux := web.NewWebMux()

	fmt.Printf("Started server on port %v", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal/api"
)

func main() {
	mux := api.NewApiMux()

	fmt.Println("Started server on port 8060")

	http.ListenAndServe(":8060", mux)
}

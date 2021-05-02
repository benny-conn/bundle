package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/repo"
)

func main() {

	port := os.Getenv("REPO_PORT")
	mux := repo.NewRepositoryMux()

	fmt.Printf("Started Repo Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

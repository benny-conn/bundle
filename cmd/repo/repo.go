package main

import (
	"fmt"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/repo"
)

func main() {

	port := os.Getenv("REPO_PORT")

	srv := repo.NewRepoServer()
	addr := fmt.Sprintf(":%v", port)
	internal.RunInternalServer(srv, addr, "repo")
}

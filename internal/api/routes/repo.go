package routes

import (
	"fmt"
	"net/http"
	"os"
)

func RepoPluginsHandlerFunc(w http.ResponseWriter, req *http.Request) {
	port := os.Getenv("REPO_PORT")
	http.Redirect(w, req, fmt.Sprintf(":%v/plugins", port), http.StatusTemporaryRedirect)
}
func RepoReadmesHandlerFunc(w http.ResponseWriter, req *http.Request) {
	port := os.Getenv("REPO_PORT")
	http.Redirect(w, req, fmt.Sprintf(":%v/readmes", port), http.StatusTemporaryRedirect)
}
func RepoThumbnailsHandlerFunc(w http.ResponseWriter, req *http.Request) {
	port := os.Getenv("REPO_PORT")
	http.Redirect(w, req, fmt.Sprintf(":%v/thumbnails", port), http.StatusTemporaryRedirect)
}

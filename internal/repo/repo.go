package repo

import (
	"net/http"

	auth "github.com/bennycio/bundle/internal/auth/user"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	readmesHandler := http.HandlerFunc(readmesHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", auth.AuthUpload(pluginsHandler))
	mux.Handle("/repo/readmes", auth.AuthUpload(readmesHandler))
	mux.Handle("/repo/thumbnails", auth.AuthUpload(thumbnailsHandler))

	return mux
}

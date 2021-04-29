package repository

import (
	"net/http"

	auth "github.com/bennycio/bundle/internal/auth/user"
	"github.com/bennycio/bundle/internal/repository/routes"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	readmesHandler := http.HandlerFunc(routes.ReadmesHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(routes.ThumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", auth.AuthUpload(pluginsHandler))
	mux.Handle("/repo/readmes", auth.AuthUpload(readmesHandler))
	mux.Handle("/repo/thumbnails", auth.AuthUpload(thumbnailsHandler))

	return mux
}

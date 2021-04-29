package repository

import (
	"net/http"

	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/internal/repository/routes"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	readmesHandler := http.HandlerFunc(routes.ReadmesHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(routes.ThumbnailsHandlerFunc)

	mux.Handle("/plugins", auth.AuthUpload(pluginsHandler))
	mux.Handle("/readmes", auth.AuthUpload(readmesHandler))
	mux.Handle("/thumbnails", auth.AuthUpload(thumbnailsHandler))

	return mux
}

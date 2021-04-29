package repo

import (
	"net/http"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	readmesHandler := http.HandlerFunc(readmesHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", authUpload(pluginsHandler))
	mux.Handle("/repo/readmes", authUpload(readmesHandler))
	mux.Handle("/repo/thumbnails", authUpload(thumbnailsHandler))

	return mux
}

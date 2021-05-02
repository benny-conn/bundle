package repo

import (
	"net/http"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", simpleAuth(pluginsHandler))
	mux.Handle("/repo/thumbnails", simpleAuth(thumbnailsHandler))

	return mux
}

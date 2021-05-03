package repo

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
)

func NewRepoServer() *http.Server {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", simpleAuth(pluginsHandler))
	mux.Handle("/repo/thumbnails", simpleAuth(thumbnailsHandler))

	return internal.MakeServerFromMux(mux)
}

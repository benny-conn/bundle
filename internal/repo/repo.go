package repo

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
)

func NewRepoServer() *http.Server {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", pluginsHandler)
	mux.Handle("/repo/thumbnails", thumbnailsHandler)

	return internal.MakeServerFromMux(mux)
}

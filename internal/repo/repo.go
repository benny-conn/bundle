package repo

import (
	"net/http"

	"github.com/bennycio/bundle/internal/api"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailsHandler := http.HandlerFunc(thumbnailsHandlerFunc)

	mux.Handle("/repo/plugins", api.SimpleAuth(authUpload(pluginsHandler)))
	mux.Handle("/repo/thumbnails", api.SimpleAuth(authUpload(thumbnailsHandler)))

	return mux
}

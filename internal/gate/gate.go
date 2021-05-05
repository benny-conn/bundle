package gate

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
)

func NewGateServer() *http.Server {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	usersHandler := http.HandlerFunc(usersHandlerFunc)
	repoPluginsHandler := http.HandlerFunc(repoPluginsHandlerFunc)
	repoThumbnailsHandler := http.HandlerFunc(repoThumbnailsHandlerFunc)
	readmesHandler := http.HandlerFunc(readmesHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", simpleAuth(usersHandler))
	mux.Handle("/api/readmes", readmesHandler)
	mux.Handle("/api/repo/plugins", authUpload(repoPluginsHandler))
	mux.Handle("/api/repo/thumbnails", authUpload(repoThumbnailsHandler))

	return internal.MakeServerFromMux(mux)
}

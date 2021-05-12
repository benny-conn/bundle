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
	sessionsHandler := http.HandlerFunc(sessionHandlerFunc)
	bundlesHandler := http.HandlerFunc(bundleHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", simpleAuth(usersHandler))
	mux.Handle("/api/readmes", basicAuth(readmesHandler))
	mux.Handle("/api/sessions", simpleAuth(sessionsHandler))
	mux.Handle("/api/bundles", simpleAuth(bundlesHandler))
	mux.Handle("/api/repo/plugins", basicAuth(repoPluginsHandler))
	mux.Handle("/api/repo/thumbnails", simpleAuth(repoThumbnailsHandler))

	return internal.MakeServerFromMux(mux)
}

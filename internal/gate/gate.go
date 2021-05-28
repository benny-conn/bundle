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
	changelogsHandler := http.HandlerFunc(changelogHandlerFunc)
	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/changelogs", changelogsHandler)
	mux.Handle("/api/users", scopedAuth(usersHandler, "users"))
	mux.Handle("/api/readmes", basicAuth(readmesHandler))
	mux.Handle("/api/sessions", scopedAuth(sessionsHandler, "sessions"))
	mux.Handle("/api/repo/plugins", basicAuth(repoPluginsHandler))
	mux.Handle("/api/repo/thumbnails", scopedAuth(repoThumbnailsHandler, "thumbnails"))

	return internal.MakeServerFromMux(mux)
}

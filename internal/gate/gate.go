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

	checkoutCompleteHandler := http.HandlerFunc(checkoutCompleteHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/purchases/complete", checkoutCompleteHandler)
	mux.Handle("/api/changelogs", basicAuth(changelogsHandler, http.MethodPost))
	mux.Handle("/api/users", scopedAuth(usersHandler, "users"))
	mux.Handle("/api/readmes", basicAuth(readmesHandler, http.MethodPost, http.MethodPatch))
	mux.Handle("/api/sessions", scopedAuth(sessionsHandler, "sessions"))
	mux.Handle("/api/repo/plugins", basicAuth(repoPluginsHandler, http.MethodPost))
	mux.Handle("/api/repo/thumbnails", scopedAuth(repoThumbnailsHandler, "thumbnails"))

	return internal.MakeServerFromMux(mux)
}

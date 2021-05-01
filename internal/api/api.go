package api

import (
	"net/http"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	usersHandler := http.HandlerFunc(usersHandlerFunc)
	repoPluginsHandler := http.HandlerFunc(repoPluginsHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", SimpleAuth(usersHandler))
	mux.Handle("/api/repo/plugins", repoPluginsHandler)

	return mux
}

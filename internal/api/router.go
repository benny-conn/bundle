package api

import (
	"net/http"

	"github.com/bennycio/bundle/internal/api/routes"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	usersHandler := http.HandlerFunc(routes.UsersHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", usersHandler)

	return mux
}

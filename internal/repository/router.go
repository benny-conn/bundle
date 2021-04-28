package repository

import (
	"net/http"

	"github.com/bennycio/bundle/internal/repository/routes"
)

func NewRepositoryMux() http.Handler {
	mux := http.NewServeMux()
	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)

	mux.Handle("/plugins", pluginsHandler)

	return mux
}

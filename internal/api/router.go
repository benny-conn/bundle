package api

import (
	"net/http"

	"github.com/bennycio/bundle/internal/api/routes"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	repoPluginsHandler := http.HandlerFunc(routes.RepoPluginsHandlerFunc)
	repoReadmesHandler := http.HandlerFunc(routes.RepoReadmesHandlerFunc)
	repoThumbnailsHandler := http.HandlerFunc(routes.RepoThumbnailsHandlerFunc)
	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	usersHandler := http.HandlerFunc(routes.UsersHandlerFunc)
	authHandler := http.HandlerFunc(routes.AuthUserHandlerFunc)

	mux.Handle("/api/repo/plugins", repoPluginsHandler)
	mux.Handle("/api/repo/readmes", repoReadmesHandler)
	mux.Handle("/api/repo/thumbnails", repoThumbnailsHandler)
	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", usersHandler)
	mux.Handle("/api/auth", authHandler)

	return mux
}

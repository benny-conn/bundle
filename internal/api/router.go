package api

import (
	"net/http"

	"github.com/bennycio/bundle/internal/api/routes"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/rs/cors"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	bundlesHandler := http.HandlerFunc(routes.BundleHandlerFunc)
	pluginsHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	usersHandler := http.HandlerFunc(routes.UsersHandlerFunc)

	mux.Handle("/api/bundles", auth.AuthUpload(bundlesHandler))
	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", usersHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "https://bundlemc.io/"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	return handler
}

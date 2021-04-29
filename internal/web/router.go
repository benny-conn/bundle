package web

import (
	"net/http"

	"github.com/bennycio/bundle/internal/web/routes"
	"github.com/rs/cors"
)

func NewWebMux() http.Handler {

	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(routes.RootHandlerFunc)
	signupHandler := http.HandlerFunc(routes.SignupHandlerFunc)
	loginHandler := http.HandlerFunc(routes.LoginHandlerFunc)
	logoutHandler := http.HandlerFunc(routes.LogoutHandlerFunc)
	pluginHandler := http.HandlerFunc(routes.PluginsHandlerFunc)
	profileHandler := http.HandlerFunc(routes.ProfileHandlerFunc)

	mux.Handle("/", rootHandler)
	mux.Handle("/plugin", pluginHandler)
	mux.Handle("/profile", profileHandler)
	mux.Handle("/login", loginHandler)
	mux.Handle("/logout", logoutHandler)
	mux.Handle("/signup", signupHandler)
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("assets/public"))))

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

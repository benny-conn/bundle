package web

import (
	"net/http"

	"github.com/rs/cors"
)

func NewWebMux() http.Handler {

	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(rootHandlerFunc)
	signupHandler := http.HandlerFunc(signupHandlerFunc)
	loginHandler := http.HandlerFunc(loginHandlerFunc)
	logoutHandler := http.HandlerFunc(logoutHandlerFunc)
	pluginHandler := http.HandlerFunc(pluginsHandlerFunc)
	profileHandler := http.HandlerFunc(profileHandlerFunc)

	mux.Handle("/", noGate(rootHandler))
	mux.Handle("/plugin", noGate(pluginHandler))
	mux.Handle("/profile", loginGate(profileHandler))
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

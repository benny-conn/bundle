package web

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
	"github.com/rs/cors"
)

func NewWebServer() *http.Server {

	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(rootHandlerFunc)
	signupHandler := http.HandlerFunc(signupHandlerFunc)
	loginHandler := http.HandlerFunc(loginHandlerFunc)
	logoutHandler := http.HandlerFunc(logoutHandlerFunc)
	pluginHandler := http.HandlerFunc(pluginsHandlerFunc)
	profileHandler := http.HandlerFunc(profileHandlerFunc)

	mux.Handle("/", rootHandler)
	mux.Handle("/plugins", pluginHandler)
	mux.Handle("/profile", loginGate(profileHandler))
	mux.Handle("/login", loginHandler)
	mux.Handle("/logout", logoutHandler)
	mux.Handle("/signup", signupHandler)
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("assets/public"))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"localhost:8080", "bundlemc.io"},
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

	handler := refresh(c.Handler(mux))

	return internal.MakeServerFromMux(handler)
}

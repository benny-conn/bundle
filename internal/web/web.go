package web

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
	"github.com/rs/cors"
)

// TODO custom routing

func NewWebServer() *http.Server {

	mux := http.NewServeMux()
	rootHandler := http.HandlerFunc(rootHandlerFunc)
	signupHandler := http.HandlerFunc(signupHandlerFunc)
	loginHandler := http.HandlerFunc(loginHandlerFunc)
	logoutHandler := http.HandlerFunc(logoutHandlerFunc)
	pluginHandler := http.HandlerFunc(pluginsHandlerFunc)
	thumbnailHandler := http.HandlerFunc(thumbnailHandlerFunc)
	bundlerHandler := http.HandlerFunc(bundlerHandlerFunc)
	profileHandler := http.HandlerFunc(profileHandlerFunc)
	ftpHandler := http.HandlerFunc(ftpHandlerFunc)
	stripeAuthHandler := http.HandlerFunc(stripeAuthHandlerFunc)
	stripeReturnHandler := http.HandlerFunc(stripeReturnHandlerFunc)

	mux.Handle("/", rootHandler)
	mux.Handle("/plugins", pluginHandler)
	mux.Handle("/plugins/", pluginHandler)
	mux.Handle("/plugins/thumbnails", thumbnailHandler)
	mux.Handle("/plugins/bundler", bundlerHandler)
	mux.Handle("/profile", loginGate(profileHandler))
	mux.Handle("/ftp", loginGate(ftpHandler))
	mux.Handle("/stripe/auth", stripeAuthHandler)
	mux.Handle("/stripe/return", stripeReturnHandler)
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

	handler := c.Handler(mux)

	return internal.MakeServerFromMux(handler)
}

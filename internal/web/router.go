package web

import (
	"net/http"

	"github.com/bennycio/bundle/internal/auth"
	"github.com/rs/cors"
)

func NewBundleMux() http.Handler {

	mux := http.NewServeMux()

	bundleHandler := http.HandlerFunc(BundleHandlerFunc)
	userHandler := http.HandlerFunc(UserHandlerFunc)
	rootHandler := http.HandlerFunc(RootHandlerFunc)
	signupHandler := http.HandlerFunc(SignupHandlerFunc)
	pluginsHandler := http.HandlerFunc(PluginsHandlerFunc)
	loginHandler := http.HandlerFunc(LoginHandlerFunc)
	logoutHandler := http.HandlerFunc(LogoutHandlerFunc)
	pluginHandler := http.HandlerFunc(PluginHandlerFunc)

	mux.Handle("/bundle", auth.AuthUpload(bundleHandler))
	mux.Handle("/users", auth.AuthWithoutScope(userHandler))
	mux.Handle("/", auth.AuthWithoutScope(rootHandler))
	mux.Handle("/signup", auth.AuthWithoutScope(signupHandler))
	mux.Handle("/plugins", auth.AuthWithoutScope(pluginsHandler))
	mux.Handle("/plugin", auth.AuthWithoutScope(pluginHandler))
	mux.Handle("/login", loginHandler)
	mux.Handle("/logout", logoutHandler)
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

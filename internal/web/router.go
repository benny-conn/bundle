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
	profileHandler := http.HandlerFunc(ProfileHandlerFunc)

	mux.Handle("/bundle", auth.AuthUpload(bundleHandler))
	mux.Handle("/users", auth.RefreshOrContinue(userHandler))
	mux.Handle("/", auth.RefreshOrContinue(rootHandler))
	mux.Handle("/plugins", pluginsHandler)
	mux.Handle("/plugin", auth.RefreshOrContinue(pluginHandler))
	mux.Handle("/profile", auth.AuthReq(profileHandler))
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

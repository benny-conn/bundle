package main

import (
	"fmt"
	"net/http"

	"github.com/bennycio/bundle/internal/auth"
	"github.com/bennycio/bundle/internal/web"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/bundle/")
	viper.AddConfigPath("$HOME/.bundle")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	mux := http.NewServeMux()

	bundleHandler := http.HandlerFunc(web.BundleHandlerFunc)
	userHandler := http.HandlerFunc(web.UserHandlerFunc)
	mux.Handle("/bundle", auth.ImplicitLogin(bundleHandler))
	mux.Handle("/users", userHandler)
	mux.HandleFunc("/", web.RootHandlerFunc)
	mux.HandleFunc("/signup", web.SignupHandlerFunc)
	mux.HandleFunc("/plugins", web.PluginHandlerFunc)
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

	fmt.Println("Started server on port 8080")

	http.ListenAndServe(":8080", handler)
}

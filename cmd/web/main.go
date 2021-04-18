package main

import (
	"fmt"
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
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
	mux.Handle("/bundle", bundle.AuthUser(bundleHandler))
	mux.Handle("/users", bundle.AuthClient(userHandler))
	mux.HandleFunc("/", web.RootHandlerFunc)
	mux.HandleFunc("/signup", web.SignupHandlerFunc)
	mux.HandleFunc("/plugins", web.PluginHandlerFunc)
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("assets/public"))))

	fmt.Println("Started server on port 8080")
	http.ListenAndServe(":8080", mux)
}

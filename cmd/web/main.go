package main

import (
	"fmt"
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
)

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

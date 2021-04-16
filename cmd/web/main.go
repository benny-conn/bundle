package web

import (
	"fmt"
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
)

func main() {
	fmt.Println("Started server on port 8080")
	mux := http.NewServeMux()

	bundleHandler := http.HandlerFunc(web.BundleHandlerFunc)
	userHandler := http.HandlerFunc(web.UserHandlerFunc)

	fmt.Println("Started server on port 8070")
	mux.Handle("/bundle", bundle.AuthUser(bundleHandler))
	mux.Handle("/users", bundle.AuthClient(userHandler))
	mux.HandleFunc("/", web.HandleRoot)
	mux.HandleFunc("/signup", web.HandleSignup)
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("../../assets/public"))))

	http.ListenAndServe(":8080", mux)
}

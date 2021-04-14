package api

import (
	"fmt"
	"net/http"

	bundle "github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/api"
)

func main() {

	mux := http.NewServeMux()

	bundleHandler := http.HandlerFunc(api.BundleHandlerFunc)
	userHandler := http.HandlerFunc(api.UserHandlerFunc)

	fmt.Println("Started server on port 8070")
	mux.Handle("/bundle", bundle.AuthUser(bundleHandler))
	mux.Handle("/users", bundle.AuthClient(userHandler))

	http.ListenAndServe(":8070", mux)
}

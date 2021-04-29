package api

import (
	"net/http"

	auth "github.com/bennycio/bundle/internal/auth/client"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	usersHandler := http.HandlerFunc(usersHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", auth.AuthClient(usersHandler))

	return mux
}

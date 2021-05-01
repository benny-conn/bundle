package api

import (
	"net/http"
)

func NewApiMux() http.Handler {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	usersHandler := http.HandlerFunc(usersHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", simpleAuth(usersHandler))

	return mux
}

package gate

import (
	"net/http"

	"github.com/bennycio/bundle/internal"
)

func NewGateMux() *http.Server {
	mux := http.NewServeMux()

	pluginsHandler := http.HandlerFunc(pluginsHandlerFunc)
	usersHandler := http.HandlerFunc(usersHandlerFunc)
	repoPluginsHandler := http.HandlerFunc(repoPluginsHandlerFunc)
	readmesHandler := http.HandlerFunc(readmesHandlerFunc)

	mux.Handle("/api/plugins", pluginsHandler)
	mux.Handle("/api/users", simpleAuth(usersHandler))
	mux.Handle("/api/readmes", readmesHandler)
	mux.Handle("/api/repo/plugins", authUpload(repoPluginsHandler))

	return internal.MakeServerFromMux(mux)
}

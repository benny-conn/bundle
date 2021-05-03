package gate

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/bennycio/bundle/internal"
)

func NewGateServer() *http.Server {
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

func newGateHttpClient() http.Client {

	clientCert, _ := tls.LoadX509KeyPair("tls/server-cert.pem", "tls/server-key.pem")
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{clientCert},
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := http.Client{
		Transport: &transport,
		Timeout:   20 * time.Second,
	}
	return client
}

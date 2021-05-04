package gate

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
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

	cert, err := ioutil.ReadFile("out/Bundle.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	clientCert, err := tls.LoadX509KeyPair("out/client.crt", "out/client.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{
		Transport: transport,
		Timeout:   1 * time.Minute,
	}
	return client
}

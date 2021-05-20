package internal

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Set timeouts so that a slow or malicious client doesn't
// hold resources forever
func MakeServerFromMux(mux http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func RunPublicServer(srv *http.Server, addr string, service string) {
	srv.Addr = addr
	if os.Getenv("MODE") == "PROD" {

		dataDir := "./tls/"

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("bundlemc.io", "*.bundlemc.io"),
			Cache:      autocert.DirCache(dataDir),
		}

		srv.TLSConfig = m.TLSConfig()

		go func() {
			err := srv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
		if service == "web" {
			srv.Addr = ":80"
			if m != nil {
				srv.Handler = m.HTTPHandler(srv.Handler)
			}
			srv.Handler = m.HTTPHandler(srv.Handler)
			fmt.Printf("Started %s server on %s\n", service, addr)
			err := srv.ListenAndServe()
			if err != nil {
				log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
			}
		}
	} else {
		fmt.Printf("Started %s server on %s\n", service, addr)
		err := srv.ListenAndServeTLS("out/server.crt", "out/server.key")
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	}
}

func RunInternalServer(srv *http.Server, addr string, service string) {

	caCertFile, err := ioutil.ReadFile("out/Bundle.crt")
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertFile)

	srv.Addr = addr
	srv.TLSConfig = &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
		MinVersion: tls.VersionTLS12,
	}

	fmt.Printf("Started %s server on %s\n", service, addr)

	err = srv.ListenAndServeTLS("out/server.crt", "out/server.key")
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}

}

func NewTlsClient() http.Client {

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

func NewBasicClient() http.Client {

	client := http.Client{
		Timeout: 1 * time.Minute,
	}
	return client
}

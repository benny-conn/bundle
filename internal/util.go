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

func MakeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func RunWebServer(srv *http.Server, addr string, service string) {

	if os.Getenv("MODE") == "PROD" {

		dataDir := "./tls/" + service

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("bundlemc.io", "*.bundlemc.io"),
			Cache:      autocert.DirCache(dataDir),
		}
		srv.Addr = ":443"
		srv.TLSConfig = m.TLSConfig()

		go func() {
			err := srv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
		if m != nil {
			srv.Handler = m.HTTPHandler(srv.Handler)
		}
		srv.Addr = addr
		srv.Handler = m.HTTPHandler(srv.Handler)
		fmt.Printf("Started %s server on %s\n", service, addr)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	} else {
		srv.Addr = addr
		fmt.Printf("Started %s server on %s\n", service, addr)
		err := srv.ListenAndServe()
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

func GetScheme() string {
	// mode := os.Getenv("MODE")

	// if mode == "PROD" {
	// 	return "https://"
	// } else {
	// 	return "http://"
	// }
	return "https://"

}

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/web"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	var srv *http.Server
	var m *autocert.Manager
	mode := os.Getenv("MODE")
	port := os.Getenv("WEB_PORT")
	if mode == "PROD" {
		srv = web.NewWebServer()
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "bundlemc.io"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		dataDir := "./tls/web"
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}
		srv.Addr = ":443"
		srv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		go func() {
			fmt.Printf("Starting HTTPS server on %s\n", srv.Addr)
			err := srv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}
	srv = web.NewWebServer()
	if m != nil {
		srv.Handler = m.HTTPHandler(srv.Handler)
	}

	srv.Addr = fmt.Sprintf(":%v", port)
	fmt.Printf("Starting HTTP server on %s\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}

}

package main

import (
	"fmt"
	"os"

	"github.com/bennycio/bundle/internal/gate"
)

func main() {
	port := os.Getenv("GATE_PORT")
	mux := gate.NewGateMux()

	fmt.Printf("Started Api Server on port %v\n", port)

	mux.Addr = fmt.Sprintf(":%v", port)
	mux.ListenAndServe()
	// mux.ListenAndServeTLS("bundlemc.io/cert.pem", "bundlemc.io/key.pem")
}

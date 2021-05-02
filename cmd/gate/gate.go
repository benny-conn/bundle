package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/gate"
)

func main() {
	port := os.Getenv("GATE_PORT")
	mux := gate.NewGateMux()

	fmt.Printf("Started Api Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

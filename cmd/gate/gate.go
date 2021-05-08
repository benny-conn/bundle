package main

import (
	"fmt"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/gate"
)

func main() {
	port := os.Getenv("GATE_PORT")
	srv := gate.NewGateServer()
	addr := fmt.Sprintf(":%v", port)
	internal.RunPublicServer(srv, addr, "gate")
}

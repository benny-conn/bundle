package main

import (
	"fmt"
	"os"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/web"
)

func main() {
	srv := web.NewWebServer()

	port := os.Getenv("WEB_PORT")

	addr := fmt.Sprintf(":%v", port)
	internal.RunPublicServer(srv, addr, "web")

}

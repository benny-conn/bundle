package main

import (
	"log"

	"github.com/bennycio/bundle/internal/mem"
)

func main() {

	err := mem.RunServer()
	if err != nil {
		log.Fatal("could not start server: " + err.Error())
	}
}

package main

import (
	"log"

	"github.com/bennycio/bundle/internal/db"
)

func main() {

	err := db.RunServer()
	if err != nil {
		log.Fatal("could not start server: " + err.Error())
	}
}

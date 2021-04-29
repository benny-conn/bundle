package main

import (
	"fmt"
	"log"

	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/db"
	"github.com/spf13/viper"
)

func init() {
	internal.InitEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/storage")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	err := db.RunServer()
	if err != nil {
		log.Fatal("could not start server: " + err.Error())
	}
}

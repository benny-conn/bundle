package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bennycio/bundle/internal/gate"
	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/gate")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	port := os.Getenv("GATE_PORT")
	mux := gate.NewGateMux()

	fmt.Printf("Started Api Server on port %v\n", port)

	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

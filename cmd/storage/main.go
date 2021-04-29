package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	internal.InitEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/storage")
	viper.AddConfigPath("/etc/bundle/")
	viper.AddConfigPath("$HOME/.bundle")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error parsing config file: %s", err))
	}
}

func main() {
	port := os.Getenv("DATABASE_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("tls/server-cert.pem", "tls/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterUsersServiceServer(grpcServer, storage.NewUsersServer())
	api.RegisterPluginsServiceServer(grpcServer, storage.NewPluginsServer())

	fmt.Printf("Started Database Server on port %v\n", port)

	grpcServer.Serve(lis)
}

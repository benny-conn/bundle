package main

import (
	"fmt"
	"log"
	"net"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/storage"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	internal.InitConfig()
}

func main() {
	port := viper.GetInt("Port")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterUsersServiceServer(grpcServer, storage.NewUsersServer())
	api.RegisterPluginsServiceServer(grpcServer, storage.NewPluginsServer())
	grpcServer.Serve(lis)
}

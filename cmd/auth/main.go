package main

import (
	"fmt"
	"log"
	"net"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	internal.InitConfig()
}
func main() {
	port := viper.GetInt("Port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterAuthServiceServer(grpcServer, auth.NewAuthServer())
	grpcServer.Serve(lis)
}

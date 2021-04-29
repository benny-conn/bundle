package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	internal.InitConfig()
}
func main() {
	port := os.Getenv("AUTH_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("tls/server-cert.pem", "tls/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterAuthServiceServer(grpcServer, auth.NewAuthServer())
	grpcServer.Serve(lis)
}

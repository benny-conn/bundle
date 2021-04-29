package db

import (
	"fmt"
	"net"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunServer() error {
	port := os.Getenv("DATABASE_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return err
	}
	creds, err := credentials.NewServerTLSFromFile("tls/server-cert.pem", "tls/server-key.pem")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterUsersServiceServer(grpcServer, newUsersServer())
	api.RegisterPluginsServiceServer(grpcServer, newPluginsServer())
	fmt.Printf("Started Database Server on port %v\n", port)

	grpcServer.Serve(lis)
	return nil
}

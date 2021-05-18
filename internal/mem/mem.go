package mem

import (
	"fmt"
	"net"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunServer() error {
	port := os.Getenv("MEM_PORT")
	host := os.Getenv("MEM_HOST")
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return err
	}

	creds, err := credentials.NewServerTLSFromFile("out/grpc/service.pem", "out/grpc/service.key")

	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterSessionServiceServer(grpcServer, newSessionsServer())
	fmt.Printf("Started Memory Storage Server on port %v\n", port)

	grpcServer.Serve(lis)
	return nil
}
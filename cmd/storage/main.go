package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/storage"
	"google.golang.org/grpc"
)

var port = 8000

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServiceServer(grpcServer, storage.NewUsersServer())
	pb.RegisterPluginsServiceServer(grpcServer, storage.NewPluginsServer())
	grpcServer.Serve(lis)
}

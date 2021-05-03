package db

import (
	"fmt"
	"net"
	"os"

	"github.com/bennycio/bundle/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Plugin struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LastUpdated int64              `bson:"lastUpdated" json:"lastUpdated"`
	api.Plugin
}

func RunServer() error {
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return err
	}
	creds, err := credentials.NewServerTLSFromFile("bundlemc.io/cert.pem", "bundlemc.io/key.pem")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	api.RegisterUsersServiceServer(grpcServer, newUsersServer())
	api.RegisterPluginsServiceServer(grpcServer, newPluginsServer())
	api.RegisterReadmeServiceServer(grpcServer, newReadmesServer())
	fmt.Printf("Started Database Server on port %v\n", port)

	grpcServer.Serve(lis)
	return nil
}

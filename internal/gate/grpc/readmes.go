package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type readmesGrpcClient interface {
	Get(req *api.Plugin) (*api.Readme, error)
	Update(req *api.Readme) error
	Insert(req *api.Readme) error
}

type readmesGrpcClientImpl struct {
	Host string
	Port string
}

func NewReadmeClient(host string, port string) readmesGrpcClient {
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
	}
	if port == "" {
		port = os.Getenv("DATABASE_PORT")
	}
	return &readmesGrpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (r *readmesGrpcClientImpl) Get(req *api.Plugin) (*api.Readme, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewReadmeServiceClient(conn)
	rdme, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return rdme, nil
}

func (r *readmesGrpcClientImpl) Update(req *api.Readme) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewReadmeServiceClient(conn)
	_, err = client.Update(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (r *readmesGrpcClientImpl) Insert(req *api.Readme) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%v:%v", r.Host, r.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewReadmeServiceClient(conn)
	_, err = client.Insert(context.Background(), req)
	if err != nil {
		return err
	}
	return nil

}

package client

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type bundlesRpcClient interface {
	Get(req *api.Bundle) (*api.Bundle, error)
	Insert(req *api.Bundle) error
	Update(req *api.Bundle) error
	Delete(req *api.Bundle) error
}

type bundlesRpcClientImpl struct {
	Host string
	Port string
}

func NewBundlesClient(host string, port string) bundlesRpcClient {
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
	}
	if port == "" {
		port = os.Getenv("DATABASE_PORT")
	}
	return &bundlesRpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (r *bundlesRpcClientImpl) Get(req *api.Bundle) (*api.Bundle, error) {
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
	client := api.NewBundleServiceClient(conn)
	ses, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return ses, nil
}

func (r *bundlesRpcClientImpl) Insert(req *api.Bundle) error {

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
	client := api.NewBundleServiceClient(conn)
	_, err = client.Insert(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (r *bundlesRpcClientImpl) Delete(req *api.Bundle) error {
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
	client := api.NewBundleServiceClient(conn)
	_, err = client.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	return nil

}

func (r *bundlesRpcClientImpl) Update(req *api.Bundle) error {
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
	client := api.NewBundleServiceClient(conn)
	_, err = client.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	return nil

}

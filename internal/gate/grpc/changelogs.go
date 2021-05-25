package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type changelogsRpcClient interface {
	Get(req *api.Changelog) (*api.Changelog, error)
	Insert(req *api.Changelog) error
}

type changelogsRpcClientImpl struct {
	Host string
	Port string
}

func NewChangelogsClient(host string, port string) changelogsRpcClient {
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
	}
	if port == "" {
		port = os.Getenv("DATABASE_PORT")
	}
	return &changelogsRpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (r *changelogsRpcClientImpl) Get(req *api.Changelog) (*api.Changelog, error) {
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
	client := api.NewChangelogServiceClient(conn)
	ses, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return ses, nil
}

func (r *changelogsRpcClientImpl) Insert(req *api.Changelog) error {

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
	client := api.NewChangelogServiceClient(conn)
	_, err = client.Insert(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

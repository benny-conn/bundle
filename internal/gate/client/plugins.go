package client

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type pluginsGrpcClient interface {
	Get(req *api.Plugin) (*api.Plugin, error)
	Update(req *api.Plugin) error
	Insert(req *api.Plugin) error
	Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error)
}

type pluginsGrpcClientImpl struct {
	Host string
	Port string
}

func NewPluginClient(host string, port string) pluginsGrpcClient {
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
	}
	if port == "" {
		port = os.Getenv("DATABASE_PORT")
	}
	return &pluginsGrpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (p *pluginsGrpcClientImpl) Get(req *api.Plugin) (*api.Plugin, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%v:%v", p.Host, p.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	pl, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return pl, nil
}
func (p *pluginsGrpcClientImpl) Insert(plugin *api.Plugin) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%v:%v", p.Host, p.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	_, err = client.Insert(context.Background(), plugin)
	if err != nil {
		return err
	}
	return nil
}
func (p *pluginsGrpcClientImpl) Update(req *api.Plugin) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%v:%v", p.Host, p.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)

	_, err = client.Update(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
func (p *pluginsGrpcClientImpl) Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%v:%v", p.Host, p.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)

	results, err := client.Paginate(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return results, nil
}

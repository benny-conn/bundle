package api

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type userGrpcService struct{}

type pluginGrpcService struct{}

func newUserClient() *userGrpcService     { return &userGrpcService{} }
func newPluginClient() *pluginGrpcService { return &pluginGrpcService{} }
func (u *userGrpcService) Get(req *api.GetUserRequest) (*api.User, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	user, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userGrpcService) Update(req *api.UpdateUserRequest) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.Update(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

func (u *userGrpcService) Insert(user *api.User) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.Insert(context.Background(), user)
	if err != nil {
		return err
	}
	return nil

}

func (p *pluginGrpcService) Get(req *api.GetPluginRequest) (*api.Plugin, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
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
func (p *pluginGrpcService) Insert(plugin *api.Plugin) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
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
func (p *pluginGrpcService) Update(req *api.UpdatePluginRequest) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
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
func (p *pluginGrpcService) Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	host := os.Getenv("DATABASE_HOST")
	addr := fmt.Sprintf("%v:%v", host, port)
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

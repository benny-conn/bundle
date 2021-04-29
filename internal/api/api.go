package api

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"google.golang.org/grpc"
)

type userGrpcService struct{}

type pluginGrpcService struct{}

func newUserClient() internal.UserService     { return &userGrpcService{} }
func newPluginClient() internal.PluginService { return &pluginGrpcService{} }

func (u *userGrpcService) Get(req *api.GetUserRequest) (*api.User, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	user, err := client.GetUser(context.Background(), req)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.UpdateUser(context.Background(), req)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewUsersServiceClient(conn)
	_, err = client.InsertUser(context.Background(), user)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	pl, err := client.GetPlugin(context.Background(), req)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)
	_, err = client.InsertPlugin(context.Background(), plugin)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)

	_, err = client.UpdatePlugin(context.Background(), req)
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
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewPluginsServiceClient(conn)

	results, err := client.PaginatePlugins(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return results, nil
}

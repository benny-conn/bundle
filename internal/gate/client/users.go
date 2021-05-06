package client

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type usersGrpcClient interface {
	Get(req *api.User) (*api.User, error)
	Update(req *api.User) error
	Insert(req *api.User) error
}

type usersGrpcClientImpl struct {
	Host string
	Port string
}

func NewUserClient(host string, port string) usersGrpcClient {
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
	}
	if port == "" {
		port = os.Getenv("DATABASE_PORT")
	}
	return &usersGrpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (u *usersGrpcClientImpl) Get(req *api.User) (*api.User, error) {
	creds, err := getCert()
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%v:%v", u.Host, u.Port)
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

func (u *usersGrpcClientImpl) Update(req *api.User) error {
	creds, err := getCert()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%v:%v", u.Host, u.Port)
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

func (u *usersGrpcClientImpl) Insert(user *api.User) error {

	creds, err := getCert()
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%v:%v", u.Host, u.Port)
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

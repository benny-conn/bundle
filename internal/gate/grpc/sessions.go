package grpc

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

type sessionsGrpcClient interface {
	Get(req *api.Session) (*api.Session, error)
	Insert(req *api.Session) (*api.SessionInsertResponse, error)
	Delete(req *api.Session) error
}

type sessionsGrpcClientImpl struct {
	Host string
	Port string
}

func NewSessionsClient(host string, port string) sessionsGrpcClient {
	if host == "" {
		host = os.Getenv("MEM_HOST")
	}
	if port == "" {
		port = os.Getenv("MEM_PORT")
	}
	return &sessionsGrpcClientImpl{
		Host: host,
		Port: port,
	}
}

func (r *sessionsGrpcClientImpl) Get(req *api.Session) (*api.Session, error) {
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
	client := api.NewSessionServiceClient(conn)
	ses, err := client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return ses, nil
}

func (r *sessionsGrpcClientImpl) Insert(req *api.Session) (*api.SessionInsertResponse, error) {
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
	client := api.NewSessionServiceClient(conn)
	res, err := client.Insert(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *sessionsGrpcClientImpl) Delete(req *api.Session) error {
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
	client := api.NewSessionServiceClient(conn)
	_, err = client.Delete(context.Background(), req)
	if err != nil {
		return err
	}
	return nil

}

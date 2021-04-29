package wrapper

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"google.golang.org/grpc"
)

func NewJWT(user *api.User) (*api.Jwt, error) {
	creds, err := GetCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("AUTH_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewAuthServiceClient(conn)
	jwt, err := client.NewJwt(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

func Refresh(jwt *api.Jwt) (*api.Jwt, error) {
	creds, err := GetCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("AUTH_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewAuthServiceClient(conn)
	ref, err := client.Refresh(context.Background(), jwt)
	if err != nil {
		return nil, err
	}
	return ref, nil
}
func GetUserFromToken(jwt *api.Jwt) (*api.User, error) {
	creds, err := GetCert()
	if err != nil {
		return nil, err
	}
	port := os.Getenv("AUTH_PORT")
	addr := fmt.Sprintf(":%v", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewAuthServiceClient(conn)
	user, err := client.GetUser(context.Background(), jwt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

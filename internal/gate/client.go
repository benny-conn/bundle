package gate

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
type pluginsGrpcClient interface {
	Get(req *api.Plugin) (*api.Plugin, error)
	Update(req *api.Plugin) error
	Insert(req *api.Plugin) error
	Paginate(req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error)
}

type readmesGrpcClient interface {
	Get(req *api.Plugin) (*api.Readme, error)
	Update(req *api.Readme) error
	Insert(req *api.Readme) error
}

type usersGrpcClientImpl struct {
	Host string
	Port string
}

type pluginsGrpcClientImpl struct {
	Host string
	Port string
}

type readmesGrpcClientImpl struct {
	Host string
	Port string
}

func newUserClient(host string, port string) usersGrpcClient {
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
func newPluginClient(host string, port string) pluginsGrpcClient {
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

func newReadmeClient(host string, port string) readmesGrpcClient {
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
		fmt.Println(err.Error())
		return err
	}

	return nil

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

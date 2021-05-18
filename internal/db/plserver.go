package db

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
)

type pluginsServer struct {
	orm *orm.PluginsOrm
	api.UnimplementedPluginsServiceServer
}

func (s *pluginsServer) Get(ctx context.Context, req *api.Plugin) (*api.Plugin, error) {

	pl, err := s.orm.Get(req)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *pluginsServer) Update(ctx context.Context, req *api.Plugin) (*api.Empty, error) {
	err := s.orm.Update(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) Insert(ctx context.Context, plugin *api.Plugin) (*api.Empty, error) {
	err := s.orm.Insert(plugin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) Paginate(ctx context.Context, req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	pls, err := s.orm.Paginate(req)
	if err != nil {

		return nil, err
	}
	return &api.PaginatePluginsResponse{
		Plugins: pls,
	}, nil
}

func newPluginsServer() *pluginsServer {
	s := &pluginsServer{orm: orm.NewPluginsOrm()}
	return s
}

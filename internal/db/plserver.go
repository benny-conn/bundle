package db

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal"
	"github.com/bennycio/bundle/internal/db/orm"
)

type pluginsServer struct {
	orm internal.PluginService
	api.UnimplementedPluginsServiceServer
}

func (s *pluginsServer) GetPlugin(ctx context.Context, req *api.GetPluginRequest) (*api.Plugin, error) {

	pl, err := s.orm.Get(req)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *pluginsServer) UpdatePlugin(ctx context.Context, req *api.UpdatePluginRequest) (*api.Empty, error) {
	err := s.orm.Update(req)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) InsertPlugin(ctx context.Context, plugin *api.Plugin) (*api.Empty, error) {
	err := s.orm.Insert(plugin)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) PaginatePlugins(ctx context.Context, req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	pls, err := s.orm.Paginate(req)
	if err != nil {
		return nil, err
	}
	return pls, nil
}

func newPluginsServer() *pluginsServer {
	s := &pluginsServer{orm: orm.NewPluginsOrm()}
	return s
}

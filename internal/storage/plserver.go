package storage

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/storage/orm"
)

type pluginsServer struct {
	api.UnimplementedPluginsServiceServer
}

func (s *pluginsServer) GetPlugin(ctx context.Context, req *api.GetPluginRequest) (*api.Plugin, error) {

	pl, err := orm.GetPlugin(req.Name)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *pluginsServer) UpdatePlugin(ctx context.Context, req *api.UpdatePluginRequest) (*api.Empty, error) {
	err := orm.UpdatePlugin(req.Name, req.UpdatedPlugin)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) InsertPlugin(ctx context.Context, plugin *api.Plugin) (*api.Empty, error) {
	err := orm.InsertPlugin(plugin)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *pluginsServer) PaginatePlugins(ctx context.Context, req *api.PaginatePluginsRequest) (*api.PaginatePluginsResponse, error) {
	pls, err := orm.PaginatePlugins(int(req.Page))
	if err != nil {
		return nil, err
	}
	return &api.PaginatePluginsResponse{
		Plugins: pls,
	}, nil
}

func NewPluginsServer() *pluginsServer {
	s := &pluginsServer{}
	return s
}

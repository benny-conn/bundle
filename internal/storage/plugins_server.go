package storage

import (
	"context"

	pb "github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/storage/orm"
)

type pluginsServer struct {
	pb.UnimplementedPluginsServiceServer
}

func (s *pluginsServer) GetPlugin(ctx context.Context, req *pb.GetPluginRequest) (*pb.Plugin, error) {

	pl, err := orm.GetPlugin(req.Name)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *pluginsServer) UpdatePlugin(ctx context.Context, req *pb.UpdatePluginRequest) (*pb.SuccessResponse, error) {
	err := orm.UpdatePlugin(req.Name, req.UpdatedPlugin)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, err
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *pluginsServer) InsertPlugin(ctx context.Context, plugin *pb.Plugin) (*pb.SuccessResponse, error) {
	err := orm.InsertPlugin(plugin)
	if err != nil {
		return &pb.SuccessResponse{Success: false}, err
	}
	return &pb.SuccessResponse{Success: true}, nil
}

func (s *pluginsServer) PaginatePlugins(ctx context.Context, req *pb.PaginatePluginsRequest) (*pb.PaginatePluginsResponse, error) {
	pls, err := orm.PaginatePlugins(int(req.Page))
	if err != nil {
		return nil, err
	}
	return &pb.PaginatePluginsResponse{
		Plugins: pls,
	}, nil
}

func NewPluginsServer() *pluginsServer {
	s := &pluginsServer{}
	return s
}

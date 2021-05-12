package db

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
)

type bundlesServer struct {
	orm *orm.BundlesOrm
	api.UnimplementedSessionServiceServer
}

func (s *bundlesServer) Get(ctx context.Context, req *api.Bundle) (*api.Bundle, error) {

	pl, err := s.orm.Get(req)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *bundlesServer) Insert(ctx context.Context, req *api.Bundle) (*api.Empty, error) {
	err := s.orm.Insert(req)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *bundlesServer) Delete(ctx context.Context, req *api.Bundle) (*api.Empty, error) {
	err := s.orm.Delete(req)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *bundlesServer) Update(ctx context.Context, req *api.Bundle) (*api.Empty, error) {
	err := s.orm.Update(req)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func newBundleServer() *bundlesServer {
	s := &bundlesServer{orm: orm.NewBundlesOrm()}
	return s
}

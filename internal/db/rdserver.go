package db

import (
	"context"
	"fmt"
	"os"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
)

type readmesServer struct {
	orm *orm.ReadmesOrm
	api.UnimplementedReadmeServiceServer
}

func (s *readmesServer) Get(ctx context.Context, req *api.Plugin) (*api.Readme, error) {

	pl, err := s.orm.Get(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "GET README"+err.Error())
		return nil, err
	}

	return pl, nil

}

func (s *readmesServer) Update(ctx context.Context, req *api.Readme) (*api.Empty, error) {
	err := s.orm.Update(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func (s *readmesServer) Insert(ctx context.Context, readme *api.Readme) (*api.Empty, error) {
	err := s.orm.Insert(readme)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func newReadmesServer() *readmesServer {
	s := &readmesServer{orm: orm.NewReadmesOrm()}
	return s
}

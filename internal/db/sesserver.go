package db

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
)

type sessionsServer struct {
	orm *orm.SessionsOrm
	api.UnimplementedSessionServiceServer
}

func (s *sessionsServer) Get(ctx context.Context, req *api.Session) (*api.Session, error) {

	pl, err := s.orm.Get(req)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func (s *sessionsServer) Insert(ctx context.Context, req *api.Session) (*api.Empty, error) {
	err := s.orm.Insert(req)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *sessionsServer) Delete(ctx context.Context, readme *api.Session) (*api.Empty, error) {
	err := s.orm.Delete(readme)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func newSessionsServer() *sessionsServer {
	s := &sessionsServer{orm: orm.NewSessionsOrm()}
	return s
}

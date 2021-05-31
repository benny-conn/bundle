package db

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
)

type changelogServer struct {
	orm *orm.ChangelogOrm
	api.UnimplementedChangelogServiceServer
}

func (s *changelogServer) Get(ctx context.Context, req *api.Changelog) (*api.Changelog, error) {

	pl, err := s.orm.Get(req)
	if err != nil {

		return nil, err
	}

	return pl, nil

}

func (s *changelogServer) Insert(ctx context.Context, req *api.Changelog) (*api.Empty, error) {
	err := s.orm.Insert(req)
	if err != nil {
		return nil, err
	}
	return &api.Empty{}, nil
}

func (s *changelogServer) GetAll(ctx context.Context, req *api.Changelog) (*api.Changelogs, error) {

	pl, err := s.orm.GetAll(req)
	if err != nil {
		return nil, err
	}

	return pl, nil

}

func newChangelogServer() *changelogServer {
	s := &changelogServer{orm: orm.NewChangelogOrm()}
	return s
}

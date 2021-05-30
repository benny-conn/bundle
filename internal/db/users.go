package db

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/db/orm"
	"github.com/bennycio/bundle/internal/logger"
)

type usersServer struct {
	orm *orm.UsersOrm
	api.UnimplementedUsersServiceServer
}

func (s *usersServer) Get(ctx context.Context, req *api.User) (*api.User, error) {
	user, err := s.orm.Get(req)
	if err != nil {

		return nil, err
	}
	return user, nil
}
func (s *usersServer) Update(ctx context.Context, req *api.User) (*api.Empty, error) {
	err := s.orm.Update(req)
	if err != nil {
		logger.ErrLog.Println(err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func (s *usersServer) Insert(ctx context.Context, user *api.User) (*api.Empty, error) {
	err := s.orm.Insert(user)
	if err != nil {
		logger.ErrLog.Println(err.Error())
		return nil, err
	}
	return &api.Empty{}, nil
}

func newUsersServer() *usersServer {
	s := &usersServer{orm: orm.NewUsersOrm()}
	return s
}

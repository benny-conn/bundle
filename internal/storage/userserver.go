package storage

import (
	"context"

	"github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/storage/orm"
)

type usersServer struct {
	api.UnimplementedUsersServiceServer
}

func (s *usersServer) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {

	user, err := orm.GetUser(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *usersServer) UpdateUser(ctx context.Context, req *api.UpdateUserRequest) (*api.Empty, error) {
	err := orm.UpdateUser(req.Username, req.UpdatedUser)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func (s *usersServer) InsertUser(ctx context.Context, user *api.User) (*api.Empty, error) {
	err := orm.InsertUser(user)
	if err != nil {
		return &api.Empty{}, err
	}
	return &api.Empty{}, nil
}

func NewUsersServer() *usersServer {
	s := &usersServer{}
	return s
}

package storage

import (
	"context"

	pb "github.com/bennycio/bundle/api"
	"github.com/bennycio/bundle/internal/storage/orm"
)

type usersServer struct {
	pb.UnimplementedUsersServiceServer
}

func (s *usersServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {

	user, err := orm.GetUser(req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *usersServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.SuccessResponse, error) {
	err := orm.UpdateUser(req.Username, req.UpdatedUser)
	if err != nil {
		return &pb.SuccessResponse{
			Success: false,
		}, err
	}
	return &pb.SuccessResponse{
		Success: true,
	}, nil
}

func (s *usersServer) InsertUser(ctx context.Context, user *pb.User) (*pb.SuccessResponse, error) {
	err := orm.InsertUser(user)
	if err != nil {
		return &pb.SuccessResponse{
			Success: false,
		}, err
	}
	return &pb.SuccessResponse{
		Success: true,
	}, nil
}

func NewUsersServer() *usersServer {
	s := &usersServer{}
	return s
}

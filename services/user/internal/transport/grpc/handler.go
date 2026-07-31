package grpc

import (
	"context"

	"github.com/disdreamq/fantastic-telegram/services/user/internal/service"
	pb "github.com/disdreamq/fantastic-telegram/services/user/proto"
	"google.golang.org/grpc"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	uSVC *service.UserService
}

func (s *userServer) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error) {
	u, err := s.uSVC.GetByID(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{Id: u.ID, Username: u.Username, Exists: true}, nil
}

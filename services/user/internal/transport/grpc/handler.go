package grpc

import (
	"context"

	"github.com/disdreamq/fantastic-telegram/services/user/internal/port"
	pb "github.com/disdreamq/fantastic-telegram/services/user/proto"
	"google.golang.org/grpc"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	uSVC port.UserService
}

func (s *userServer) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error) {
	ctx = context.WithValue(ctx, "trace_id", in.Request_ID)
	u, err := s.uSVC.GetByID(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{ID: u.ID, Username: u.Username, Exists: true, Request_ID: in.Request_ID}, nil
}

package grpc

import (
	"context"
	"strings"

	"github.com/disdreamq/fantastic-telegram/services/user/internal/port"
	pb "github.com/disdreamq/fantastic-telegram/services/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
	p port.TokenProvider
}

func (s *userServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	authValues := strings.Split(in.Token, " ")
	if len(authValues) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing auth header")
	}
	if len(authValues) != 2 || authValues[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "missing auth header")
	}
	token := authValues[1]
	payload, err := s.p.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return &pb.ValidateTokenResponse{Valid: true, UserID: payload.Claims.UserID, UserEmail: payload.Claims.Email}, nil
}

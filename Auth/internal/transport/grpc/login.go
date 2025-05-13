package transport

import (
	myerrors "auth/internal/error"
	"auth/pkg/proto"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	tokens, err := s.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to login user")
	}

	return &proto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

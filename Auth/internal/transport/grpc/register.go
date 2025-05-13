package transport

import (
	myerrors "auth/internal/error"
	"auth/pkg/proto"
	"context"
	"errors"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (s *serverAPI) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username can not be empty")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email can not be empty")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password can not be empty")
	}

	if !emailRegex.MatchString(req.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email format")
	}

	user, err := s.auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, myerrors.ErrUserAlreadyExists.Error())
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &proto.RegisterResponse{
		UserId:   user.UUID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

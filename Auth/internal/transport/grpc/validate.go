package transport

import (
	"auth/pkg/logger"
	"auth/pkg/proto"
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	if len(req.AccessToken) == 0 || !(strings.HasPrefix(req.AccessToken, "Bearer ")) {
		return nil, status.Error(codes.InvalidArgument, "invalid jwt token")
	}

	accessTokenString := req.AccessToken[7:] // извлечение самой строки токена

	user, err := s.auth.VerifyToken(ctx, accessTokenString) 
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to validate user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to validate user")
	}

	return &proto.VerifyTokenResponse{
		UserId:   user.UserID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

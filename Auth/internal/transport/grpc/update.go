package transport

import (
	"auth/pkg/proto"
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) UpdateTokens(ctx context.Context, req *proto.UpdateTokensRequest) (*proto.UpdateTokensResponse, error) {
	if len(req.RefreshToken) == 0 || !(strings.HasPrefix(req.RefreshToken, "Bearer ")) {
		return nil, status.Error(codes.InvalidArgument, "invalid jwt token")
	}

	refreshTokenString := req.RefreshToken[7:] // извлечение самой строки токена

	pairToken, err := s.auth.UpdateToken(ctx, refreshTokenString)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user's token")
	}

	return &proto.UpdateTokensResponse{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken}, nil
}

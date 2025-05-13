package service

import (
	"auth/internal/entity"
	"auth/pkg/logger"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *authService) VerifyToken(ctx context.Context, accessToken string) (*entity.TokenClaims, error) {
	ClaimsJWT, err := s.checkAccessToken(accessToken)
	if err != nil || ClaimsJWT == nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByUsernameOrEmail(ctx, ClaimsJWT.Username, "")
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "unable to GetUserByUsernameOrEmail", zap.Error(err))
		return nil, fmt.Errorf("unable to GetUserByUsernameOrEmail %v", err)
	}

	return &entity.TokenClaims{
		UserID:   user.UUID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

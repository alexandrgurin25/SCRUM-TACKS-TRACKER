package service

import (
	"auth/internal/entity"
	errors "auth/internal/error"

	"auth/pkg/logger"
	"context"

	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) Login(ctx context.Context, username, password string) (*entity.TokenPair, error) {
	user, err := s.userRepo.GetUserByUsernameOrEmail(ctx, username, "")
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "unable to GetUserByUsernameOrEmail", zap.Error(err))
		return nil, fmt.Errorf("unable to GetUserByUsernameOrEmail %v", err)
	}

	if user == nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed user login:", zap.Error(errors.ErrUserNotFound))
		return nil, errors.ErrUserNotFound
	}


	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed user login", zap.String("username", user.Username), zap.Error(err))
		return nil, errors.ErrIncorrectPassword
	}
	

	accessToken, err := s.generateAccessToken(user.UUID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.UUID, user.Username)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

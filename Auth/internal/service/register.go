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

func (s *authService) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) (*entity.User, error) {
	existingUser, err := s.userRepo.GetUserByUsernameOrEmail(ctx, username, email)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to check existing user", zap.Error(err))
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(ctx, username, email, string(passwordHash))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, err
}

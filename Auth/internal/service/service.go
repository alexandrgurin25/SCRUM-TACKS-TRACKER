package service

import (
	"auth/internal/config"
	"auth/internal/entity"
	repository "auth/internal/repository/postgres"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type AuthService interface {
	Register(ctx context.Context, username string, email string, password string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.TokenPair, error)
	UpdateToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error)
	VerifyToken(ctx context.Context, accessToken string) (*entity.TokenClaims, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

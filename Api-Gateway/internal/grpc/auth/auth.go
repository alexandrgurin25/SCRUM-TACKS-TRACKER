package auth

import (
	"context"
	"fmt"
	"gateway/internal/domain/models"

	auth "github.com/AlexMickh/scrum-protos/pkg/api/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	conn   *grpc.ClientConn
	client auth.AuthClient
}

func New(addr string) (*AuthService, error) {
	const op = "grpc.auth.New"

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := auth.NewAuthClient(conn)

	return &AuthService{
		conn:   conn,
		client: client,
	}, nil
}

func (a *AuthService) Register(ctx context.Context, email, username, password string) (models.User, error) {
	const op = "grpc.auth.Register"

	res, err := a.client.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{
		ID:       res.GetUserId(),
		Email:    res.GetEmail(),
		Username: res.GetUsername(),
	}, nil
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	const op = "grpc.auth.Login"

	res, err := a.client.Login(ctx, &auth.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (a *AuthService) UpdateTokens(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "grpc.auth.UpdateTokens"

	res, err := a.client.UpdateTokens(ctx, &auth.UpdateTokensRequest{RefreshToken: refreshToken})
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}

func (a *AuthService) Stop() {
	a.conn.Close()
}

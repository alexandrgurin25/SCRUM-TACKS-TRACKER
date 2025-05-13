package service

import (
	"auth/internal/entity"
	"context"
	"fmt"
)

func (s *authService) UpdateToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {

	claimsJWT, err := s.checkRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	newJwtToken, err := s.generateAccessToken(claimsJWT.UserUUID, claimsJWT.Username)
	if err != nil {
		return nil, fmt.Errorf("could not generate new token: %v", err)
	}

	newRefreshToken, err := s.generateRefreshToken(claimsJWT.UserUUID, claimsJWT.Username)
	if err != nil {
		return nil, fmt.Errorf("could not generate new refresh token: %v", err)
	}

	return &entity.TokenPair{
		AccessToken:  newJwtToken,
		RefreshToken: newRefreshToken,
	}, nil
}

package service

import (
	myerrors "auth/internal/error"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserUUID string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *authService) generateAccessToken(UserUUID string, username string) (string, error) {
	if username == "" {
		return "", myerrors.ErrFieldUserEmpty
	}
	claims := jwt.MapClaims{
		"user_id":  UserUUID,
		"username": username,
		"exp":      time.Now().Add(s.cfg.Auth.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.AccessTokenSecret))
}

func (s *authService) generateRefreshToken(UserUUID string, username string) (string, error) {
	if username == "" {
		return "", myerrors.ErrFieldUserEmpty
	}
	claims := jwt.MapClaims{
		"user_id":  UserUUID,
		"username": username,
		"exp":      time.Now().Add(s.cfg.Auth.RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.RefreshTokenSecret))
}

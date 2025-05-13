package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *authService) checkAccessToken(inToken string) (*Claims, error) {
	token, err := jwt.Parse(inToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorUnverifiable)
		}
		return []byte(s.cfg.Auth.AccessTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid access token claims")
	}

	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("access token has expired")
	}

	return &Claims{
		UserUUID: claims["user_id"].(string),
		Username: claims["username"].(string),
	}, nil
}

func (s *authService) checkRefreshToken(inToken string) (*Claims, error) {
	token, err := jwt.Parse(inToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorUnverifiable)
		}
		return []byte(s.cfg.Auth.RefreshTokenSecret), nil // Приведение к []byte
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims) // Изменено на jwt.MapClaims
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("refresh token has expired")
	}

	return &Claims{
		UserUUID: claims["user_id"].(string),
		Username: claims["username"].(string),
	}, nil
}

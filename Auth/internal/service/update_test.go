package service

import (
	"auth/internal/config"
	"auth/internal/repository/postgres/mocks"
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdate_Failed_CheckRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)

	s := NewAuthService(mockRepo, &config.Config{})

	// Задаем ожидаемые значения
	refreshToken := "invalid_refresh_token"

	_, err := s.UpdateToken(context.Background(), refreshToken)
	// Проверяем результаты
	assert.Error(t, err)
}

func TestUpdate_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)

	a := config.Auth{
		RefreshTokenSecret: "lms",
		RefreshTokenTTL:    time.Hour,
	}

	cfg := config.Config{
		Auth: a,
	}

	service := NewAuthService(mockRepo, &cfg)
	userId := uuid.New().String()
	username := "user"
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      time.Now().Add(cfg.Auth.RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Задаем ожидаемые значения
	refreshToken, _ := token.SignedString([]byte(cfg.RefreshTokenSecret))

	_, err := service.UpdateToken(context.Background(), refreshToken)
	// Проверяем результаты

	assert.NoError(t, err)
}


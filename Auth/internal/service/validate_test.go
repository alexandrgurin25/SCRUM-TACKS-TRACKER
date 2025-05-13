package service

import (
	"auth/internal/config"
	"auth/internal/entity"
	"auth/internal/repository/postgres/mocks"
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestValidate_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)

	a := config.Auth{
		AccessTokenSecret: "lms",
		AccessTokenTTL:    time.Hour,
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
		"exp":      time.Now().Add(cfg.Auth.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Задаем ожидаемые значения
	accessToken, _ := token.SignedString([]byte(cfg.AccessTokenSecret))

	mockRepo.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&entity.User{
			UUID:     userId,
			Username: username,
		}, nil,
	).AnyTimes()

	_, err := service.VerifyToken(context.Background(), accessToken)
	// Проверяем результаты
	assert.NoError(t, err)
}


func TestValidate_Failed_JwtSiging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockUserRepository(ctrl)

	a := config.Auth{
		AccessTokenSecret: "lms",
		AccessTokenTTL:    time.Hour,
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
	accessToken, _ := token.SignedString([]byte(cfg.RefreshTokenSecret))

	mockRepo.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&entity.User{
			UUID:     userId,
			Username: username,
		}, nil,
	).AnyTimes()

	_, err := service.VerifyToken(context.Background(), accessToken)
	// Проверяем результаты
	assert.Error(t, err)
}

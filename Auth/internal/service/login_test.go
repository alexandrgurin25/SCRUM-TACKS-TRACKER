package service

import (
	"auth/internal/config"
	"auth/internal/entity"
	myerrors "auth/internal/error"
	"auth/internal/repository/postgres/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Successful(t *testing.T) {
	ctx := context.Background()

	username := "username"
	email := "email@email.ru"
	password := "1234"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entity.User{
		UUID:         uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo, &config.Config{})
	mockRepo.EXPECT().
		GetUserByUsernameOrEmail(ctx, username, "").
		Return(user, nil)

	tokenPair, err := service.Login(ctx, username, password)
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}

func TestLogin_Failed_UserNotFound(t *testing.T) {
	ctx := context.Background()

	username := ""
	password := "1234"

	user := &entity.User{}
	user = nil

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo, &config.Config{})
	mockRepo.EXPECT().
		GetUserByUsernameOrEmail(ctx, username, "").
		Return(user, nil)

	tokenPair, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.ErrorIs(t, err, myerrors.ErrUserNotFound)
	assert.Nil(t, tokenPair)
}

func TestLogin_Failed_BcryptHashing(t *testing.T) {
	ctx := context.Background()

	username := ""
	password := "1234"

	user := &entity.User{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo, &config.Config{})
	mockRepo.EXPECT().
		GetUserByUsernameOrEmail(ctx, username, "").
		Return(user, nil)

	tokenPair, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

func TestLogin_Failed_GenegateAccessToken(t *testing.T) {
	ctx := context.Background()

	username := "1"
	password := "1234"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		PasswordHash: string(passwordHash),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo, &config.Config{})
	mockRepo.EXPECT().
		GetUserByUsernameOrEmail(ctx, username, "").
		Return(user, nil)

	tokenPair, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

func TestLogin_Failed_UserRepositoryError(t *testing.T) {
	ctx := context.Background()

	username := "username"

	password := "1234"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockRepo, &config.Config{})
	mockRepo.EXPECT().
		GetUserByUsernameOrEmail(ctx, username, "").
		Return(nil, errors.New("failed to get username"))

	tokenPair, err := service.Login(ctx, username, password)
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

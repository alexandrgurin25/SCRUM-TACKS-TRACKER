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

func TestRegister_Successful(t *testing.T) {
	ctx := context.Background()

	username := "test"
	email := "email"
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
		GetUserByUsernameOrEmail(ctx, username, email).
		Return(nil, nil)

	mockRepo.EXPECT().
		CreateUser(ctx, username, email, gomock.Any()).
		Return(user, nil)

	user, err := service.Register(ctx, username, email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestRegister_Failed_UserAlreadyExist(t *testing.T) {
	ctx := context.Background()

	username := "test"
	email := "email"
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
		GetUserByUsernameOrEmail(ctx, username, email).
		Return(user, nil)

	mockRepo.EXPECT().
		CreateUser(ctx, username, email, gomock.Any()).
		Return(user, nil).AnyTimes()

	user, err := service.Register(ctx, username, email, password)
	assert.Error(t, err)
	assert.ErrorIs(t, err, myerrors.ErrUserAlreadyExists)
	assert.Nil(t, user)
}

func TestRegister_Failed_GetUserByUsernameOrEmail(t *testing.T) {
	ctx := context.Background()

	username := "test"
	email := "email"
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
		GetUserByUsernameOrEmail(ctx, username, email).
		Return(nil, errors.New("falied to get user by username"))

	mockRepo.EXPECT().
		CreateUser(ctx, username, email, gomock.Any()).
		Return(user, nil).AnyTimes()

	user, err := service.Register(ctx, username, email, password)
	assert.Error(t, err)
	assert.Nil(t, user)
}


func TestRegister_Failed_CreateUser(t *testing.T) {
	ctx := context.Background()

	username := "test"
	email := "email"
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
		GetUserByUsernameOrEmail(ctx, username, email).
		Return(nil, nil)

	mockRepo.EXPECT().
		CreateUser(ctx, username, email, gomock.Any()).
		Return(user, errors.New("falied to create user"))

	user, err := service.Register(ctx, username, email, password)
	assert.Error(t, err)
	assert.Nil(t, user)
}
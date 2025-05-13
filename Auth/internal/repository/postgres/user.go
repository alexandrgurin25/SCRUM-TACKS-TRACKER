package repository

import (
	"auth/internal/entity"
	"context"

	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=user.go -destination=mocks/user.go -package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, username string, email string, passwordHash string) (*entity.User, error)
	GetUserByUsernameOrEmail(ctx context.Context, username string, email string) (*entity.User, error)
}

type DBTX interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) *userRepository {
	return &userRepository{db: db}
}

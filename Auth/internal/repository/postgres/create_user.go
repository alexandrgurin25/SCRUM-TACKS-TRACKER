package repository

import (
	"auth/internal/entity"
	"context"
	"fmt"
)

func (r *userRepository) CreateUser(ctx context.Context,
	username string,
	email string,
	passwordHash string) (*entity.User, error) {

	var user entity.User
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users (username, email, passwordHash)
		VALUES ($1, $2, $3)
		RETURNING uuid, email, username`,
		username,
		email,
		passwordHash,
	).Scan(&user.UUID, &user.Email, &user.Username)

	if err != nil {
		return nil, fmt.Errorf("unable CreateUser %v", err)
	}

	return &user, nil
}

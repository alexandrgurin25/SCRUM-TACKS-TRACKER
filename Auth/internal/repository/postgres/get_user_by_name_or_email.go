package repository

import (
	"auth/internal/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *userRepository) GetUserByUsernameOrEmail(
	ctx context.Context,
	username string,
	email string,
) (*entity.User, error) {

	var user entity.User

	err := r.db.QueryRow(
		ctx,
		`SELECT uuid, username, email, passwordHash FROM users WHERE username = $1 OR email = $2`,
		username,
		email,
	).Scan(&user.UUID, &user.Username, &user.Email, &user.PasswordHash)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("FindUserByUsername repository error -> %v", err)
	}

	return &user, nil
}

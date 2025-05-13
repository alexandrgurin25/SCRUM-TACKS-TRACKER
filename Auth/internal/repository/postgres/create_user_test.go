package repository

import (
	"auth/internal/config"
	"auth/pkg/postgres"
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Successful(t *testing.T) {
	ctx := context.Background()

	pathToCfg := "../../../config/.env"
	cfg, err := config.NewTest(pathToCfg)
	assert.NotEmpty(t, cfg.Postgres.Host, "")
	assert.NoError(t, err)

	pathToMigrate := "file://../../../db/migrations"
	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	username := "username"
	email := "user@scrum.ru"
	passwordHash := "@#$%"

	//Начинаем транзакцию
	tx, err := db.Begin(ctx)
	assert.NoError(t, err)
	defer tx.Rollback(ctx)

	repository := NewUserRepository(tx)

	user, err := repository.CreateUser(ctx, username, email, passwordHash)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	var count int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE username = $1`, username).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	if err := tx.Rollback(ctx); err != nil {
		t.Fatalf("failed to rollback transaction: %v", err)
	}
}

func TestCreateUser_Failed(t *testing.T) {
	ctx := context.Background()

	pathToCfg := "../../../config/.env"
	cfg, err := config.NewTest(pathToCfg)
	assert.NotEmpty(t, cfg.Postgres.Host, "")
	assert.NoError(t, err)

	pathToMigrate := "file://../../../db/migrations"
	db, err := postgres.NewTest(ctx, cfg, pathToMigrate)
	assert.NoError(t, err)

	username := "username"
	email := "user@scrum.ru"
	passwordHash := "@#$%"

	//Начинаем транзакцию
	tx, err := db.Begin(ctx)
	assert.NoError(t, err)
	defer tx.Rollback(ctx)

	repository := NewUserRepository(tx)

	_, err = repository.CreateUser(ctx, username, email, passwordHash)
	assert.NoError(t, err)

	_, err = repository.CreateUser(ctx, username, email, passwordHash)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable CreateUser")

	if err := tx.Rollback(ctx); err != nil {
		t.Fatalf("failed to rollback transaction: %v", err)
	}
}

package repository

import (
	"auth/internal/config"
	"auth/pkg/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByName_Successful(t *testing.T) {
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

	user, err := repository.GetUserByUsernameOrEmail(ctx, username, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	if err := tx.Rollback(ctx); err != nil {
		t.Fatalf("failed to rollback transaction: %v", err)
	}
}

func TestGetUserByName_Failed(t *testing.T) {
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

	//Начинаем транзакцию
	tx, err := db.Begin(ctx)
	assert.NoError(t, err)
	defer tx.Rollback(ctx)

	repository := NewUserRepository(tx)

	user, err := repository.GetUserByUsernameOrEmail(ctx, username, email)
	assert.NoError(t, err)
	assert.Nil(t, user)

	if err := tx.Rollback(ctx); err != nil {
		t.Fatalf("failed to rollback transaction: %v", err)
	}
}

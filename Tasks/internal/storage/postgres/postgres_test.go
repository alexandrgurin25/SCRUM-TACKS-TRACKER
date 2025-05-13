package postgres

import (
	"context"
	"task/internal/config"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	cfg := config.DBConfig{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Password:       "root",
		Name:           "tasks",
		MinPools:       3,
		MaxPools:       5,
		MigrationsPath: "../../../migrations",
	}

	postgres, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}
	defer postgres.Close()
}

func TestSaveTask(t *testing.T) {
	ctx := context.Background()
	postgres := connect(ctx, t)
	defer postgres.Close()
	id := uuid.New().String()

	_, err := postgres.SaveTask(ctx, id, "1", "1", "hi", "hello", time.Now())
	if err != nil {
		t.Fatalf("failed to seve task: %v", err)
	}
}

func TestGetTasksByBoardID(t *testing.T) {
	ctx := context.Background()
	postgres := connect(ctx, t)
	defer postgres.Close()
	boardID := "1"

	_, err := postgres.GetTasksByBoardID(ctx, boardID)
	if err != nil {
		t.Fatalf("failed to get tasks by board id: %v", err)
	}
}

func TestGetAllTasks(t *testing.T) {
	ctx := context.Background()
	postgres := connect(ctx, t)
	defer postgres.Close()

	_, err := postgres.GetAllTasks(ctx)
	if err != nil {
		t.Fatalf("failed to get all tasks: %v", err)
	}
}

func TestChangeBoard(t *testing.T) {
	ctx := context.Background()
	postgres := connect(ctx, t)
	defer postgres.Close()
	c := make(chan error)
	id := "5"
	boardID := "2"

	go postgres.ChangeBoard(ctx, id, boardID, c)
	err := <-c
	close(c)
	if err != nil {
		t.Fatalf("failed to change board: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	ctx := context.Background()
	postgres := connect(ctx, t)
	defer postgres.Close()
	c := make(chan error)
	id := "5"

	go postgres.DeleteTask(ctx, id, c)
	err := <-c
	close(c)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}
}

func connect(ctx context.Context, t *testing.T) *Postgres {
	cfg := config.DBConfig{
		Host:           "localhost",
		Port:           5432,
		User:           "postgres",
		Password:       "root",
		Name:           "tasks",
		MinPools:       3,
		MaxPools:       5,
		MigrationsPath: "../../../migrations",
	}

	postgres, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	return postgres
}

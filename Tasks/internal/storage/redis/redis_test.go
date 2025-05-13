package redis

import (
	"context"
	"task/internal/config"
	"task/internal/domain/models"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	cfg := config.RedisConfig{
		Addr:     "localhost:6379",
		User:     "root",
		Password: "root",
		DB:       0,
	}

	redis, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to the redis: %v", err)
	}
	defer redis.Close()
}

func TestSaveTask(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()
	task := models.Task{
		ID:          "9",
		Title:       "hi",
		Description: "hello",
		AuthorID:    "1",
		BoardID:     "10",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Deadline:    time.Time{},
	}

	err := redis.SaveTask(ctx, task)
	if err != nil {
		t.Fatalf("failed to save task: %v", err)
	}
}

func TestSaveTasks(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()
	tasks := []models.Task{
		{
			ID:          "1",
			Title:       "hi",
			Description: "hello",
			AuthorID:    "1",
			BoardID:     "10",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			Deadline:    time.Time{},
		},
		{
			ID:          "2",
			Title:       "hi",
			Description: "hello",
			AuthorID:    "1",
			BoardID:     "10",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			Deadline:    time.Time{},
		},
		{
			ID:          "3",
			Title:       "hi",
			Description: "hello",
			AuthorID:    "1",
			BoardID:     "10",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
			Deadline:    time.Time{},
		},
	}

	err := redis.SaveTasks(ctx, tasks)
	if err != nil {
		t.Fatalf("failed to save tasks: %v", err)
	}
}

func TestGetTasksByBoardID(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()

	tasks, err := redis.GetTasksByBoardID(ctx, "10")
	if err != nil {
		t.Fatalf("failed to get tasks: %v", err)
	}
	t.Log(tasks)
}

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()

	task := models.Task{
		ID:          "5",
		Title:       "holla",
		Description: "sdcfalkjfkl;",
		AuthorID:    "1",
		BoardID:     "10",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Deadline:    time.Time{},
	}

	err := redis.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("failed to update tasks: %v", err)
	}
}

func TestChangeBoard(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()

	err := redis.ChangeBoard(ctx, "5", "15")
	if err != nil {
		t.Fatalf("failed to change board: %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	ctx := context.Background()
	redis := connect(ctx, t)
	defer redis.Close()

	err := redis.DeleteTask(ctx, "5")
	if err != nil {
		t.Fatalf("failed to delete tasks: %v", err)
	}
}

func connect(ctx context.Context, t *testing.T) *Redis {
	cfg := config.RedisConfig{
		Addr:     "localhost:6379",
		User:     "root",
		Password: "root",
		DB:       0,
	}
	redis, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to the redis: %v", err)
	}
	return redis
}

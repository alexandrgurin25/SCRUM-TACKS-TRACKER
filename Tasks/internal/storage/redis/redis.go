package redis

import (
	"context"
	"fmt"
	"strings"
	"task/internal/config"
	"task/internal/domain/models"
	"task/internal/storage"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	const op = "storage.redis.New"

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Redis{rdb: rdb}, nil
}

func (r *Redis) SaveTask(ctx context.Context, task models.Task) error {
	const op = "storage.redis.SaveTask"

	err := r.rdb.HSet(ctx, fmt.Sprintf("%s&%s", task.ID, task.BoardID), task).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var cursor uint64
	var keys []string
	pipe := r.rdb.Pipeline()
	for {
		keys, _, err = pipe.Scan(ctx, cursor, fmt.Sprintf("*&%s", task.BoardID), 10).Result()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if cursor == 0 {
			break
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, key := range keys {
		err = pipe.Expire(ctx, key, 24*time.Hour).Err()
		if err != nil {
			continue // ? Maybe we need queue here
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) SaveTasks(ctx context.Context, tasks []models.Task) error {
	const op = "storage.redis.SaveTasks"

	var err error
	pipe := r.rdb.Pipeline()
	for _, task := range tasks {
		err = pipe.HSet(ctx, fmt.Sprintf("%s&%s", task.ID, task.BoardID), task).Err()
		if err != nil {
			continue // ? Maybe we need queue here
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, task := range tasks {
		err = pipe.Expire(ctx, fmt.Sprintf("%s&%s", task.ID, task.BoardID), 24*time.Hour).Err()
		if err != nil {
			continue // ? Maybe we need queue here
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error) {
	const op = "storage.redis.GetTasksByBoardID"

	var cursor uint64
	var keys []string
	var err error
	pipe := r.rdb.Pipeline()
	for {
		keys, cursor, err = pipe.Scan(ctx, cursor, fmt.Sprintf("*&%s", boardID), 10).Result()
		if err != nil {
			return []models.Task{}, fmt.Errorf("%s: %w", op, err)
		}
		if cursor == 0 {
			break
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var tasks []models.Task
	for _, key := range keys {
		var task models.Task
		err = pipe.HGetAll(ctx, key).Scan(&task)
		if err != nil {
			return []models.Task{}, fmt.Errorf("%s: %w", op, err)
		}
		task.ID = strings.Split(key, "&")[0]
		task.BoardID = boardID
		tasks = append(tasks, task)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (r *Redis) UpdateTask(ctx context.Context, task models.Task) error {
	const op = "storage.redis.UpdateTask"

	err := r.rdb.HSet(ctx, fmt.Sprintf("%s&%s", task.ID, task.BoardID), task).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) ChangeBoard(ctx context.Context, id, boardID string) error {
	const op = "storage.redis.ChangeBoard"

	var cursor uint64
	keys, _, err := r.rdb.Scan(ctx, cursor, fmt.Sprintf("%s&*", id), 1).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(keys) == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	var task models.Task
	err = r.rdb.HGetAll(ctx, keys[0]).Scan(&task)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = r.rdb.Del(ctx, keys[0]).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	key := fmt.Sprintf("%s&%s", strings.Split(keys[0], "&")[0], boardID)
	err = r.rdb.HSet(ctx, key, task).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) DeleteTask(ctx context.Context, id string) error {
	const op = "storage.redis.DeleteTask"

	var cursor uint64
	keys, _, err := r.rdb.Scan(ctx, cursor, fmt.Sprintf("%s&*", id), 1).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if len(keys) == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	err = r.rdb.Del(ctx, keys[0]).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) Close() {
	r.rdb.Close()
}

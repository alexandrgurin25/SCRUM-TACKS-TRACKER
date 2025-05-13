package redis

import (
	"comments/internal/config"
	"comments/internal/domain/models"
	"context"
	"fmt"
	"strings"
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

func (r *Redis) SaveComment(ctx context.Context, comment *models.Comment) error {
	const op = "storage.redis.SaveComment"

	err := r.rdb.HSet(ctx, fmt.Sprintf("%s&%s", comment.ID, comment.TaskID), *comment).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var cursor uint64
	var keys []string
	pipe := r.rdb.Pipeline()
	for {
		keys, _, err = pipe.Scan(ctx, cursor, fmt.Sprintf("*&%s", comment.TaskID), 10).Result()
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
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) SaveComments(ctx context.Context, comments []*models.Comment) error {
	const op = "storage.redis.SaveComments"

	var err error
	pipe := r.rdb.Pipeline()
	for _, comment := range comments {
		err = pipe.HSet(ctx, fmt.Sprintf("%s&%s", comment.ID, comment.TaskID), comment).Err()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, comment := range comments {
		err = pipe.Expire(ctx, fmt.Sprintf("%s&%s", comment.ID, comment.TaskID), 24*time.Hour).Err()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

// ! this method doesn't work
func (r *Redis) GetAllByTaskId(ctx context.Context, taskId string, lastTime time.Time) ([]*models.Comment, error) {
	const op = "storage.redis.GetAllByTaskId"

	var cursor uint64
	var keys []string
	var err error
	pipe := r.rdb.Pipeline()
	for {
		keys, cursor, err = pipe.Scan(ctx, cursor, fmt.Sprintf("*&%s", taskId), 10).Result()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if cursor == 0 {
			break
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(keys) > 0 {
		fmt.Printf("keys: %v\n", keys)
	}

	var comments []*models.Comment
	for _, key := range keys {
		var comment models.Comment
		err = pipe.HGetAll(ctx, key).Scan(&comment)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if comment.CreatedAt.Unix() > lastTime.Unix() {
			continue
		}
		comment.ID = strings.Split(key, "&")[0]
		comment.TaskID = taskId
		comments = append(comments, &comment)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("get")

	return comments, nil
}

func (r *Redis) Close() {
	r.rdb.Close()
}

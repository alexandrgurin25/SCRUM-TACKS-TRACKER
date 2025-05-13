package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"task/internal/config"
	"task/internal/domain/models"
	"task/internal/storage"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.MaxPools,
		cfg.MinPools,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		strings.Split(connString, "&")[0],
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{pool: pool}, nil

}

func (p *Postgres) SaveTask(
	ctx context.Context,
	id string,
	authorID string,
	boardID string,
	title string,
	description string,
	deadline time.Time,
) (models.Task, error) {
	const op = "storage.postgres.SaveTask"

	var task models.Task
	err := p.pool.QueryRow(
		ctx,
		`INSERT INTO tasks
		(id, title, description, author_id, board_id, deadline)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, title, description, author_id, board_id, deadline, created_at, updated_at`,
		id, title, boardID, authorID, description, deadline,
	).Scan(
		&task.ID,
		&task.Title,
		&task.BoardID,
		&task.AuthorID,
		&task.Description,
		&task.Deadline,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (p *Postgres) GetTasksByBoardID(ctx context.Context, boardID string) ([]models.Task, error) {
	const op = "storage.postgres.GetTasksByBoardID"

	tasks, err := p.getAllTasks(ctx, boardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (p *Postgres) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	const op = "storage.postgres.GetAllTasks"

	tasks, err := p.getAllTasks(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (p *Postgres) UpdateTask(
	ctx context.Context,
	id string,
	title string,
	description string,
	deadline time.Time,
) (models.Task, error) {
	const op = "storage.postgres.UpdateTask"

	var task models.Task
	err := p.pool.QueryRow(
		ctx,
		`UPDATE tasks 
		SET title = $1, 
		description = $2, 
		deadline = $3, 
		updated_at = $4 
		WHERE id = $5 RETURNING *`,
		title, description, deadline, time.Now(), id,
	).Scan(
		&task.ID,
		&task.AuthorID,
		&task.BoardID,
		&task.Title,
		&task.Description,
		&task.Deadline,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (p *Postgres) ChangeBoard(ctx context.Context, id, boardID string, c chan error) {
	const op = "storage.postgres.ChangeBoard"

	_, err := p.pool.Exec(
		ctx,
		"UPDATE tasks SET board_id = $1, updated_at = $2 WHERE id = $3",
		boardID, time.Now(), id,
	)
	if err != nil {
		c <- fmt.Errorf("%s: %w", op, err)
	}

	c <- nil
}

func (p *Postgres) DeleteTask(ctx context.Context, id string, c chan error) {
	const op = "storage.postgres.DeleteTask"

	_, err := p.pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		c <- fmt.Errorf("%s: %w", op, err)
	}

	c <- nil
}

func (p *Postgres) getAllTasks(ctx context.Context, options ...string) ([]models.Task, error) {
	var rows pgx.Rows
	var err error
	if len(options) == 0 {
		rows, err = p.pool.Query(
			ctx,
			"SELECT id, title, description, author_id, board_id, deadline, created_at, updated_at FROM tasks",
		)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = p.pool.Query(
			ctx,
			`SELECT id, title, description, author_id, board_id, deadline, created_at, updated_at 
			FROM tasks 
			WHERE board_id = $1`,
			options[0],
		)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.AuthorID,
			&task.BoardID,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}
